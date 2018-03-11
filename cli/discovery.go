package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/nickw444/miio-go/common"
)

func installDiscovery(app *kingpin.Application) {
	cmd := app.Command("discover", "Discover devices on the local network")
	cmd.Action(func(ctx *kingpin.ParseContext) error {
		sharedClient.SetDiscoveryInterval(time.Second * 2)

		sub, err := sharedClient.NewSubscription()
		if err != nil {
			panic(err)
		}
		for event := range sub.Events() {
			switch event.(type) {
			case common.EventNewDevice:
				dev := event.(common.EventNewDevice).Device
				go writeDeviceInfo(dev)
			case common.EventNewMaskedDevice:
				deviceId := event.(common.EventNewMaskedDevice).DeviceID
				go writeMaskedDeviceInfo(deviceId)

			}
		}
		return nil
	})
}

func writeDeviceInfo(dev common.Device) {
	deviceInfo, _ := dev.GetInfo()
	fmt.Println("-------------")
	fmt.Println("Discovered new device:")
	fmt.Printf("ID: %d\n", dev.ID())
	fmt.Printf("Firmware Version: %s\n", deviceInfo.FirmwareVersion)
	fmt.Printf("Hardware Version: %s\n", deviceInfo.HardwareVersion)
	fmt.Printf("Mac Address: %s\n", deviceInfo.MacAddress)
	fmt.Printf("Model: %s\n", deviceInfo.Model)
	fmt.Printf("Token: %s\n", hex.EncodeToString(dev.GetToken()))
	fmt.Println("-------------")
}

func writeMaskedDeviceInfo(deviceId uint32) {
	fmt.Println("-------------")
	fmt.Println("Discovered new device with masked token:")
	fmt.Printf("ID: %d\n", deviceId)
	fmt.Println("You must manually retreive this token in order to communicate with the device.")
	fmt.Println("-------------")
}
