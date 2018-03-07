package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/nickw444/miio-go/protocol/packet"
	"github.com/nickw444/miio-go/protocol/transport"
	"github.com/nickw444/miio-go/simulator/capability"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type SimulatedDevice interface {
	HandlePacket(pkt *packet.Packet) (*packet.Packet, error)
	HandleDiscover(pkt *packet.Packet) (*packet.Packet, error)
}

type BaseDevice struct {
	capabilities []capability.Capability
	crypto       packet.Crypto
	deviceToken  []byte
	deviceID     uint32
	revealToken  bool
}

func NewBaseDevice(deviceID uint32, deviceToken []byte, revealToken bool) (*BaseDevice, error) {
	crypto, err := packet.NewCrypto(deviceID, deviceToken, 1, time.Now(), clock.New())
	if err != nil {
		return nil, err
	}
	return &BaseDevice{
		capabilities: []capability.Capability{},
		deviceID:     deviceID,
		deviceToken:  deviceToken,
		crypto:       crypto,
		revealToken:  revealToken,
	}, nil
}

func (b *BaseDevice) DecodeRequest(pkt *packet.Packet) (*transport.Request, error) {
	err := b.crypto.VerifyPacket(pkt)
	if err != nil {
		panic(err)
	}

	data, err := b.crypto.Decrypt(pkt.Data)
	if err != nil {
		panic(err)
	}

	request := transport.Request{}
	err = json.Unmarshal(data, &request)
	return &request, err
}

func (b *BaseDevice) PackResponse(response interface{}) (*packet.Packet, error) {
	data, err := json.Marshal(&response)
	if err != nil {
		return nil, err
	}

	log.Infof("Response Data: %s", string(data))
	return b.crypto.NewPacket(data)
}

func (b *BaseDevice) HandleDiscover(pkt *packet.Packet) (*packet.Packet, error) {
	var checksumValue []byte
	if b.revealToken {
		checksumValue = b.deviceToken
	} else {
		checksumValue = bytes.Repeat([]byte{0x00}, 16)
	}
	return packet.New(b.deviceID, checksumValue, 1, []byte{}), nil
}

func (b *BaseDevice) getPropFromCapabilities(propName string) (interface{}, error) {
	for _, c := range b.capabilities {
		handled, result, err := c.MaybeGetProp(propName)
		if err != nil {
			return nil, err
		}

		if handled {
			return result, nil
		}
	}

	return nil, fmt.Errorf("No capabilities available to return data for get_prop '%s'", propName)
}

func (b *BaseDevice) HandlePacket(pkt *packet.Packet) (*packet.Packet, error) {
	req, err := b.DecodeRequest(pkt)
	if err != nil {
		return nil, err
	}

	log.Infof("Request received. ID=%d method=%s, params=%s", req.ID, req.Method, req.Params)

	switch req.Method {
	case "get_prop":
		props := req.Params.([]interface{})
		retProps := []interface{}{}

		for _, prop := range props {
			propName := prop.(string)
			value, err := b.getPropFromCapabilities(propName)
			if err != nil {
				return nil, err
			}

			retProps = append(retProps, value)
		}

		return b.PackResponse(transport.Response{
			ID:     req.ID,
			Result: retProps,
		})

	default:
		for _, c := range b.capabilities {
			handled, result, err := c.MaybeHandle(req.Method, req.Params)
			if err != nil {
				return nil, err
			}
			if handled {
				return b.PackResponse(transport.Response{
					ID:     req.ID,
					Result: result,
				})
			}
		}

		log.Warnf("No capabilities able to handle method %s", req.Method)
	}

	return nil, nil
}

func (b *BaseDevice) AddCapability(c capability.Capability) {
	b.capabilities = append(b.capabilities, c)
}
