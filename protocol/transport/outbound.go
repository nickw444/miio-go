package transport

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"sync"

	"github.com/benbjohnson/clock"
	"github.com/nickw444/miio-go/common"
	"github.com/nickw444/miio-go/protocol/packet"
)

type OutboundConn interface {
	WriteTo([]byte, net.Addr) (int, error)
}

// Outbound transport is an abstraction around a net.UDPConn for outbound interaction with
// a networked miIO device. Consumers should never close the underlying socket and continue
// to use the service. Outbound also provides retry and timeout logic.
type Outbound interface {
	// Handle handles incoming packets and triggers waiting continuations.
	Handle(pkt *packet.Packet) error
	// Call makes a call, waits for a response and returns the raw bytes returned.
	Call(method string, params interface{}) ([]byte, error)
	// CallAndDeserialize makes a call, waits for a response and deserialises the JSON
	// payload into `ret`.
	CallAndDeserialize(method string, params interface{}, resp interface{}) error
	// Send will send a raw packet without waiting for a response.
	Send(packet *packet.Packet) error
}

type requestID uint32

type outbound struct {
	maxRetries int
	timeout    time.Duration

	clock  clock.Clock
	crypto packet.Crypto

	dest   net.Addr
	socket OutboundConn

	nextReqID          requestID
	continuationsMutex sync.RWMutex
	continuations      map[requestID]chan []byte
}

func NewOutbound(crypto packet.Crypto, dest net.Addr, socket OutboundConn) Outbound {
	return newOutbound(10, time.Millisecond*200, clock.New(), crypto, dest, socket)
}

func newOutbound(maxRetries int, timeout time.Duration, clock clock.Clock, crypto packet.Crypto,
	dest net.Addr, socket OutboundConn) *outbound {
	return &outbound{
		maxRetries: maxRetries,
		timeout:    timeout,
		clock:      clock,
		crypto:     crypto,
		dest:       dest,
		socket:     socket,

		nextReqID:     1,
		continuations: make(map[requestID]chan []byte),
	}
}

func (o *outbound) Handle(pkt *packet.Packet) error {
	if pkt.Header.Length <= 32 {
		return nil
	}

	err := o.crypto.VerifyPacket(pkt)
	if err != nil {
		panic(err)
	}

	data, err := o.crypto.Decrypt(pkt.Data)
	if err != nil {
		panic(err)
	}

	resp := Response{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}

	// Lookup the response ID and pass data to the appropriate continuation goroutine.
	o.continuationsMutex.RLock()
	if ch, ok := o.continuations[resp.ID]; ok {
		common.Log.Debugf("Callback with ID %d was reconciled", resp.ID)
		ch <- data
	} else {
		common.Log.Debugf("Unable to reconcile callback for resp id %d", resp.ID)
	}
	o.continuationsMutex.RUnlock()

	return nil
}

func (o *outbound) Call(method string, params interface{}) ([]byte, error) {
	defer func() { o.nextReqID++ }()

	// Setup a continuation channel
	o.continuationsMutex.Lock()
	ch := make(chan []byte)
	o.continuations[o.nextReqID] = ch
	o.continuationsMutex.Unlock()

	// Ensure we cleanup.
	defer func() {
		o.continuationsMutex.Lock()
		delete(o.continuations, o.nextReqID)
		close(ch)
		o.continuationsMutex.Unlock()
	}()

	for i := 0; i < o.maxRetries+1; i++ {
		// Perform the call
		err := o.call(o.nextReqID, method, params)
		if err != nil {
			return nil, err
		}

		select {
		case data := <-ch:
			return data, nil
		case <-o.clock.After(o.timeout):
			common.Log.Debugf("Timed out whilst waiting for response.")
			continue
		}
	}

	err := fmt.Errorf("Max retries exceeded whilst sending request to device %s", o.dest)
	common.Log.Error(err)
	return nil, err
}

func (o *outbound) CallAndDeserialize(method string, params interface{}, ret interface{}) error {
	resp, err := o.Call(method, params)
	err = json.Unmarshal(resp, ret)
	if err != nil {
		return err
	}
	return nil
}

func (o *outbound) Send(packet *packet.Packet) error {
	common.Log.Debugf("Sending packet with checksum: %s", hex.EncodeToString(packet.Header.Checksum))
	_, err := o.socket.WriteTo(packet.Serialize(), o.dest)
	return err
}

// Call out to the device, but don't wait for a response.
func (o *outbound) call(requestId requestID, method string, params interface{}) (err error) {
	data, err := json.Marshal(Request{
		ID:     requestId,
		Method: method,
		Params: params,
	})
	if err != nil {
		return
	}

	p, err := o.crypto.NewPacket(data)
	if err != nil {
		return
	}

	err = o.Send(p)
	return
}

type Response struct {
	ID     requestID   `json:"id"`
	Result interface{} `json:"result"`
}

type Request struct {
	ID     requestID   `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}
