package capability

import (
	"strconv"

	"github.com/nickw444/miio-go/common"
	"github.com/nickw444/miio-go/protocol/transport"
	"github.com/nickw444/miio-go/subscription"
)

type Light struct {
	subscriptionTarget subscription.SubscriptionTarget
	outbound           transport.Outbound

	lastEvent common.EventUpdateLight
}

func NewLight(target subscription.SubscriptionTarget, transport transport.Outbound) *Light {
	return &Light{
		subscriptionTarget: target,
		outbound:           transport,
	}
}

func (l *Light) SetBrightness(brightness int) error {
	_, err := l.outbound.Call("set_bright", []interface{}{brightness})
	if err != nil {
		return err
	}

	l.lastEvent.Brightness = brightness
	return l.subscriptionTarget.Publish(l.lastEvent)
}

func (l *Light) SetHSV(hue int, saturation int) error {
	_, err := l.outbound.Call("set_hsv", []interface{}{hue, saturation})
	if err != nil {
		return err
	}

	l.lastEvent.Hue = hue
	return l.subscriptionTarget.Publish(l.lastEvent)
}

func (l *Light) SetRGB(red int, green int, blue int) error {
	rgb := miioRGB(0)
	rgb.SetComponents(red, green, blue)
	_, err := l.outbound.Call("set_rgb", []interface{}{int(rgb)})
	if err != nil {
		return err
	}

	l.lastEvent.RGB.Red = red
	l.lastEvent.RGB.Green = green
	l.lastEvent.RGB.Blue = blue
	return l.subscriptionTarget.Publish(l.lastEvent)
}

func (l *Light) Update() error {
	var resp transport.Response
	props := []string{"bright", "color_mode", "rgb", "hue", "sat"}
	err := l.outbound.CallAndDeserialize("get_prop", props, &resp)
	if err != nil {
		return err
	}

	didUpdate := false
	for i, result := range resp.Result.([]interface{}) {
		propName := props[i]
		result, _ := strconv.Atoi(result.(string))
		switch propName {
		case "bright":
			if l.lastEvent.Brightness != result {
				didUpdate = true
				l.lastEvent.Brightness = result
			}
		case "color_mode":
			if l.lastEvent.ColorMode != result {
				didUpdate = true
				l.lastEvent.ColorMode = result
			}
		case "rgb":
			rgb := miioRGB(result)
			red, green, blue := rgb.GetComponents()
			if l.lastEvent.RGB.Red != red || l.lastEvent.RGB.Green != green || l.lastEvent.RGB.Blue != blue {
				didUpdate = true
				l.lastEvent.RGB.Red = red
				l.lastEvent.RGB.Green = green
				l.lastEvent.RGB.Blue = blue
			}
		case "hue":
			if l.lastEvent.Hue != result {
				didUpdate = true
				l.lastEvent.Hue = result
			}
		case "sat":
			if l.lastEvent.Saturation != result {
				didUpdate = true
				l.lastEvent.Saturation = result
			}
		}
	}
	if didUpdate {
		return l.subscriptionTarget.Publish(l.lastEvent)
	}
	return nil
}

type miioRGB int

func (m *miioRGB) GetComponents() (red int, green int, blue int) {
	red = int(*m) >> 16 & 0xff
	green = int(*m) >> 8 & 0xff
	blue = int(*m) & 0xff
	return
}

func (m *miioRGB) SetComponents(red int, green int, blue int) {
	i := 0
	i |= red << 16
	i |= green << 8
	i |= blue
	*m = miioRGB(i)
}
