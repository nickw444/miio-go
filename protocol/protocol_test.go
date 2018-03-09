package protocol

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/nickw444/miio-go/device"
	deviceMocks "github.com/nickw444/miio-go/device/mocks"
	"github.com/nickw444/miio-go/protocol/packet"
	"github.com/nickw444/miio-go/protocol/transport"
	transportMocks "github.com/nickw444/miio-go/protocol/transport/mocks"
	subscriptionMocks "github.com/nickw444/miio-go/subscription/common/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Protocol_SetUp() (tt struct {
	clk                *clock.Mock
	transport          *mockTransport
	deviceFactory      DeviceFactory
	cryptoFactory      CryptoFactory
	subscriptionTarget *subscriptionMocks.SubscriptionTarget
	protocol           *protocol
	devices            []*deviceMocks.Device
	broadcastDevice    *deviceMocks.Device
}) {
	tt.clk = clock.NewMock()
	tt.transport = &mockTransport{new(transportMocks.Inbound)}
	tt.subscriptionTarget = new(subscriptionMocks.SubscriptionTarget)
	tt.deviceFactory = func(deviceId uint32, outbound transport.Outbound, seen time.Time) device.Device {
		d := &deviceMocks.Device{}
		tt.devices = append(tt.devices, d)
		return d
	}
	tt.cryptoFactory = func(deviceID uint32, deviceToken []byte, initialStamp uint32, stampTime time.Time) (packet.Crypto, error) {
		return nil, nil
	}
	tt.broadcastDevice = &deviceMocks.Device{}
	tt.broadcastDevice.On("Discover").Return(nil)
	tt.protocol = newProtocol(tt.clk, tt.transport, tt.deviceFactory, tt.cryptoFactory, tt.subscriptionTarget, tt.broadcastDevice)
	return
}

// Ensure that the broadcast device has Discover called on it.
func TestProtocol_Discover(t *testing.T) {
	tt := Protocol_SetUp()

	err := tt.protocol.Discover()
	assert.NoError(t, err)
	tt.broadcastDevice.AssertCalled(t, "Discover")
}

// Ensure that inbound's Packets method is called.
func TestProtocol_dispatcher(t *testing.T) {
	tt := Protocol_SetUp()
	wg := sync.WaitGroup{}
	wg.Add(1)

	ch := make(chan *packet.Packet)
	// Hack to convert the channel to a read-only channel (what the mock expects)
	ro := func(c chan *packet.Packet) <-chan *packet.Packet {
		return c
	}

	tt.transport.inbound.On("Packets").Return(ro(ch)).Run(func(args mock.Arguments) {
		wg.Done()
	})
	tt.protocol.start()
	wg.Wait()
	tt.transport.inbound.AssertExpectations(t)
}

type mockTransport struct {
	inbound *transportMocks.Inbound
}

func (m *mockTransport) Inbound() transport.Inbound {
	return m.inbound
}

func (*mockTransport) NewOutbound(crypto packet.Crypto, dest net.Addr) transport.Outbound {
	return &transportMocks.Outbound{}
}

func (*mockTransport) Close() error {
	return nil
}
