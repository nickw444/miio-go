package main

import (
	"bytes"
	"encoding/hex"
	"net"

	"github.com/alecthomas/kingpin"
	"github.com/nickw444/miio-go/protocol/packet"
	"github.com/nickw444/miio-go/protocol/transport"
	"github.com/nickw444/miio-go/simulator/device"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	defaultToken := bytes.Repeat([]byte{0x00, 0xff}, 8)
	var (
		deviceType  = kingpin.Arg("device", "Device to simulate").Default("yeelight").Enum("yeelight", "powerplug")
		deviceId    = kingpin.Flag("device-id", "Device ID for the simulated device").Default("12341234").Uint32()
		deviceToken = kingpin.Flag("device-token", "The device token to use for encrypted payloads").Default(hex.EncodeToString(defaultToken)).HexBytes()
		revealToken = kingpin.Flag("reveal-token", "Whether or not to reveal the device token").Default("true").Bool()
	)

	kingpin.Parse()

	listenAddr := &net.UDPAddr{Port: 54321}
	s, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		panic(err)
	}

	inbound := transport.NewInbound(s)
	log.Infof("Creating device with id=%d token=%s revealToken=%b",
		*deviceId, hex.EncodeToString(*deviceToken), *revealToken)
	baseDev, err := device.NewBaseDevice(*deviceId, *deviceToken, *revealToken)
	if err != nil {
		panic(err)
	}

	var dev device.SimulatedDevice
	if *deviceType == "yeelight" {
		dev = device.NewSimulatedYeelight(baseDev)
	} else if *deviceType == "powerplug" {
		dev = device.NewSimulatedPowerPlug(baseDev)
	}

	for pkt := range inbound.Packets() {
		var resp *packet.Packet
		var err error
		if pkt.Header.DeviceID == 0xffffffff {
			log.Info("Discovery packet received")
			resp, err = dev.HandleDiscover(pkt)
		} else {
			resp, err = dev.HandlePacket(pkt)
		}

		if err != nil {
			panic(err)
		}
		if resp != nil {
			s.WriteToUDP(resp.Serialize(), pkt.Meta.Addr)
		}
	}
}
