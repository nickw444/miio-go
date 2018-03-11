package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/nickw444/miio-go/capability"
	"github.com/nickw444/miio-go/common"
	"github.com/nickw444/miio-go/device"
)

var sharedDevice common.Device

func findDevice(deviceId uint32, timeout time.Duration) (common.Device, error) {
	var timeoutCh <-chan time.Time
	timeoutCh = time.After(timeout)
	sub, err := sharedClient.NewSubscription()
	if err != nil {
		panic(err)
	}
	events := sub.Events()
	defer sub.Close()

	for {
		select {
		case event := <-events:
			switch event.(type) {
			case common.EventNewDevice:
				dev := event.(common.EventNewDevice).Device
				if dev.ID() == deviceId {
					return dev, nil
				}
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("Timed out whilst connecting to device with id %d", deviceId)
		}
	}
}

func installControl(app *kingpin.Application) {
	controlCmd := app.Command("control", "Control lights")
	deviceId := controlCmd.Flag("device-id", "The ID of the device to control").Required().Uint32()

	controlCmd.Action(func(ctx *kingpin.ParseContext) (err error) {
		sharedDevice, err = findDevice(*deviceId, time.Second*5)
		return
	})

	installBrightness(controlCmd)
	installPower(controlCmd)
	installColor(controlCmd)
}

func installBrightness(parent *kingpin.CmdClause) {
	cmd := parent.Command("brightness", "Set device brightness")
	brightness := cmd.Arg("brightness", "The brightness to set (between 0-100)").Required().Int()
	cmd.Action(func(ctx *kingpin.ParseContext) error {
		var light *capability.Light

		switch sharedDevice.(type) {
		case *device.Yeelight:
			light = sharedDevice.(*device.Yeelight).Light
		default:
			return fmt.Errorf("Device with type %T cannot have brightness adjusted", sharedDevice)
		}

		return light.SetBrightness(*brightness)
	})
}

func installPower(parent *kingpin.CmdClause) {
	cmd := parent.Command("power", "Set device power")
	state := cmd.Arg("state", "The power state (on/off)").Required().Enum("on", "off")
	cmd.Action(func(ctx *kingpin.ParseContext) error {
		var power *capability.Power

		switch sharedDevice.(type) {
		case *device.Yeelight:
			power = sharedDevice.(*device.Yeelight).Power
		case *device.PowerPlug:
			power = sharedDevice.(*device.PowerPlug).Power
		default:
			return fmt.Errorf("Device with type %T cannot have brightness adjusted", sharedDevice)
		}

		return power.SetPower(common.PowerState(*state))
	})
}

func installColor(parent *kingpin.CmdClause) {
	cmd := parent.Command("color", "Set device color")

	hsv := cmd.Command("hsv", "Set color using HSV values")
	hue := hsv.Arg("hue", "Hue to set (0-360)").Required().Int()
	sat := hsv.Arg("saturation", "Saturation to set (0-100)").Required().Int()

	rgb := cmd.Command("rgb", "Set color using RGB values")
	red := rgb.Arg("red", "Red value to set (0-255)").Required().Int()
	green := rgb.Arg("green", "Green value to set (0-255)").Required().Int()
	blue := rgb.Arg("blue", "Blue value to set (0-255)").Required().Int()

	var light *capability.Light

	rgb.Action(func(ctx *kingpin.ParseContext) error {
		return light.SetRGB(*red, *green, *blue)
	})

	hsv.Action(func(ctx *kingpin.ParseContext) error {
		return light.SetHSV(*hue, *sat)
	})

	cmd.Action(func(ctx *kingpin.ParseContext) error {
		switch sharedDevice.(type) {
		case *device.Yeelight:
			light = sharedDevice.(*device.Yeelight).Light
		default:
			return fmt.Errorf("Device with type %T cannot have brightness adjusted", sharedDevice)
		}
		return nil
	})
}
