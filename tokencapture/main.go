package main

import (
	"net"

	"time"

	"flag"

	"encoding/hex"

	"github.com/nickw444/miio-go/common"
	"github.com/nickw444/miio-go/device"
	"github.com/nickw444/miio-go/protocol/packet"
	"github.com/nickw444/miio-go/protocol/tokens"
	"github.com/nickw444/miio-go/protocol/transport"
	"github.com/sirupsen/logrus"
)

var (
	log            = logrus.New()
	ignoredDevices = make(map[uint32]bool)
	t              transport.Transport
	quitChan       chan struct{}
	broadcastDev   device.Device
	tokenStore     tokens.TokenStore

	tokenStoreFile = flag.String("file", "tokens.txt", "Path to the token store to update")
)

func main() {
	miioDebug := flag.Bool("miio-debug", false, "Enable miio debug")
	flag.Parse()

	if *miioDebug {
		miioLogger := logrus.New()
		miioLogger.SetLevel(logrus.DebugLevel)
		common.SetLogger(miioLogger)
	}

	var err error
	tokenStore, err = tokens.FromFile(*tokenStoreFile)
	if err != nil {
		log.Panic(err)
	}

	var listenAddr *net.UDPAddr
	s, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		panic(err)
	}
	t = transport.NewTransport(s)

	addr := &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 54321,
	}
	devTsp := t.NewOutbound(nil, addr)
	broadcastDev = device.New(0, devTsp, time.Time{}, nil)

	go dispatcher()

	tick := time.Tick(5 * time.Second)
	broadcastDev.Discover()
	for {
		select {
		case <-quitChan:
			return
		default:
		}
		select {
		case <-quitChan:
			return
		case <-tick:
			broadcastDev.Discover()
		}
	}
}

func dispatcher() {
	pkts := t.Inbound().Packets()
	for {
		select {
		case <-quitChan:
			return
		default:
		}

		select {
		case <-quitChan:
			return
		case pkt := <-pkts:
			go process(pkt)
		}
	}
}

func process(pkt *packet.Packet) {
	if _, ok := ignoredDevices[pkt.Header.DeviceID]; ok {
		return
	}

	if pkt.DataLength() == 0 {
		if pkt.HasZeroChecksum() {
			log.Warnf("Device with Id %d is not revealing its token. Reset this device and connect to its network to retrieve the token. Ignoring it.", pkt.Header.DeviceID)
			ignoredDevices[pkt.Header.DeviceID] = true
			return
		} else {
			tokenStore.AddDevice(pkt.Header.DeviceID, pkt.Header.Checksum)
			err := tokenStore.WriteFile(*tokenStoreFile)
			if err != nil {
				log.Panic(err)
			}
			log.Infof("Got token for device with id %d: %s", pkt.Header.DeviceID, hex.EncodeToString(pkt.Header.Checksum))
		}
	}
}
