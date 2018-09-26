package protocol

import (
	"net"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/nickw444/miio-go/common"
	"github.com/nickw444/miio-go/device"
	"github.com/nickw444/miio-go/protocol/packet"
	"github.com/nickw444/miio-go/protocol/tokens"
	"github.com/nickw444/miio-go/protocol/transport"
	"github.com/nickw444/miio-go/subscription"
)

type Protocol interface {
	subscription.SubscriptionTarget

	Discover() error
	SetExpiryTime(duration time.Duration)
}

type protocol struct {
	subscription.SubscriptionTarget
	port          int
	expireAfter   time.Duration
	clock         clock.Clock
	lastDiscovery time.Time
	tokenStore    tokens.TokenStore

	broadcastDev   device.Device
	quitChan       chan struct{}
	devicesMutex   sync.RWMutex
	devices        map[uint32]device.Device
	ignoredDevices map[uint32]bool

	transport     transport.Transport
	deviceFactory DeviceFactory
	cryptoFactory CryptoFactory
}

type DeviceFactory func(deviceId uint32, outbound transport.Outbound, seen time.Time, token []byte) device.Device
type CryptoFactory func(deviceID uint32, deviceToken []byte, initialStamp uint32, stampTime time.Time) (packet.Crypto, error)

type ProtocolConfig struct {
	// Required config
	BroadcastIP net.IP
	TokenStore  tokens.TokenStore

	// Optional config
	ListenPort int // Defaults to a random system-assigned port if not provided.
}

func NewProtocol(c ProtocolConfig) (Protocol, error) {
	clk := clock.New()
	var listenAddr *net.UDPAddr
	if c.ListenPort != 0 {
		listenAddr = &net.UDPAddr{Port: c.ListenPort}
	}

	s, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	t := transport.NewTransport(s)
	deviceFactory := func(deviceId uint32, outbound transport.Outbound, seen time.Time, token []byte) device.Device {
		return device.New(deviceId, outbound, seen, token)
	}
	cryptoFactory := func(deviceID uint32, deviceToken []byte, initialStamp uint32, stampTime time.Time) (packet.Crypto, error) {
		return packet.NewCrypto(deviceID, deviceToken, initialStamp, stampTime, clk)
	}

	addr := &net.UDPAddr{
		IP:   c.BroadcastIP,
		Port: 54321,
	}
	broadcastDev := deviceFactory(0, t.NewOutbound(nil, addr), time.Time{}, nil)

	p := newProtocol(clk, t, deviceFactory, cryptoFactory, subscription.NewTarget(), broadcastDev, c.TokenStore)
	p.start()
	return p, nil
}

func newProtocol(c clock.Clock, transport transport.Transport, deviceFactory DeviceFactory,
	crptoFactory CryptoFactory, target subscription.SubscriptionTarget, broadcastDev device.Device,
	tokenStore tokens.TokenStore) *protocol {

	p := &protocol{
		SubscriptionTarget: target,
		transport:          transport,
		deviceFactory:      deviceFactory,
		cryptoFactory:      crptoFactory,
		clock:              c,
		quitChan:           make(chan struct{}),
		devices:            make(map[uint32]device.Device),
		broadcastDev:       broadcastDev,
		tokenStore:         tokenStore,
		ignoredDevices:     make(map[uint32]bool),
	}
	return p
}

func (p *protocol) start() {
	go p.dispatcher()
}

func (p *protocol) SetExpiryTime(duration time.Duration) {
	p.expireAfter = duration
}

func (p *protocol) dispatcher() {
	pkts := p.transport.Inbound().Packets()
	for {
		select {
		case <-p.quitChan:
			return
		default:
		}

		select {
		case <-p.quitChan:
			return
		case pkt := <-pkts:
			go p.process(pkt)
		}
	}
}

func (p *protocol) Discover() error {
	common.Log.Debugf("Running discovery...")

	if p.lastDiscovery.After(time.Time{}) {
		// If the device has not been seen recently, it should be expired.
		cutoff := time.Now().Add(p.expireAfter * -1)
		var expiredDevices []device.Device
		p.devicesMutex.RLock()
		for _, dev := range p.devices {
			if dev.Seen().Before(cutoff) {
				common.Log.Debugf("Device %d is stale. Last Seen at %s", dev.ID(), dev.Seen())
				expiredDevices = append(expiredDevices, dev)
			}
		}
		p.devicesMutex.RUnlock()

		for _, dev := range expiredDevices {
			common.Log.Debugf("Removing expired device with id %d.", dev.ID())
			p.removeDevice(dev.ID())
			dev.Close()
			err := p.Publish(common.EventExpiredDevice{dev})
			if err != nil {
				common.Log.Warn(err)
			}
		}
	}
	if err := p.broadcastDev.Discover(); err != nil {
		return err
	}

	p.lastDiscovery = time.Now()
	return nil
}
func (p *protocol) process(pkt *packet.Packet) {
	common.Log.Debugf("Processing incoming packet from %s", pkt.Meta.Addr)
	if ok, _ := p.ignoredDevices[pkt.Header.DeviceID]; ok {
		return
	}

	dev := p.getDevice(pkt.Header.DeviceID)
	if dev == nil && pkt.DataLength() == 0 {
		// Device response to a Hello packet.
		common.Log.Debugf("Device with id %d responded to Hello packet.", pkt.Header.DeviceID)

		deviceToken := pkt.Header.Checksum
		if pkt.HasZeroChecksum() {
			token, err := p.tokenStore.GetToken(pkt.Header.DeviceID)
			if err != nil {
				common.Log.Warnf("Device with id %d is not revealing its token. You must manually collect this token and add it to the store.", pkt.Header.DeviceID)
				p.ignoredDevices[pkt.Header.DeviceID] = true
				p.Publish(common.EventNewMaskedDevice{DeviceID: pkt.Header.DeviceID})
				return
			} else {
				common.Log.Debugf("Loaded token for device %d from store", pkt.Header.DeviceID)
				deviceToken = token
			}
		}

		crypto, err := p.cryptoFactory(pkt.Header.DeviceID, deviceToken, pkt.Header.Stamp,
			pkt.Meta.DecodeTime)
		if err != nil {
			panic(err)
		}

		t := p.transport.NewOutbound(crypto, pkt.Meta.Addr)
		baseDev := p.deviceFactory(pkt.Header.DeviceID, t, pkt.Meta.DecodeTime, deviceToken)

		// Store the provisional device for now to ensure it can handle subsequent
		// packets that may occur during classification.
		p.addDevice(baseDev)

		common.Log.Infof("Classifying device...")
		dev, err := device.Classify(baseDev)
		if err != nil {
			panic(err)
		}

		// Store the specific device and publish a new device event.
		p.addDevice(dev)
		p.Publish(common.EventNewDevice{Device: dev})
	} else if dev != nil {
		// Known device. Handle the incoming packet.
		err := dev.Handle(pkt)
		if err != nil {
			common.Log.Errorf("Unable to process packet %v for device %d. Error %s", pkt, dev.ID(), err)
		}
	} else {
		common.Log.Errorf("Unable to process packet %v. Device unknown.", pkt)
	}
}

func (p *protocol) removeDevice(id uint32) {
	p.devicesMutex.Lock()
	delete(p.devices, id)
	p.devicesMutex.Unlock()
}

func (p *protocol) addDevice(dev device.Device) {
	p.devicesMutex.Lock()
	p.devices[dev.ID()] = dev
	p.devicesMutex.Unlock()
}

func (p *protocol) getDevice(id uint32) device.Device {
	p.devicesMutex.RLock()
	dev, ok := p.devices[id]
	p.devicesMutex.RUnlock()
	if !ok {
		return nil
	}
	return dev
}
