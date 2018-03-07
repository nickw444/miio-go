package device

import (
	"github.com/nickw444/miio-go/simulator/capability"
)

type SimulatedYeelight struct {
	*BaseDevice
}

func NewSimulatedYeelight(baseDevice *BaseDevice) *SimulatedYeelight {
	baseDevice.AddCapability(&capability.Info{
		Model: "yeelink.light.color1",
	})
	baseDevice.AddCapability(&capability.Power{})

	return &SimulatedYeelight{
		BaseDevice: baseDevice,
	}
}
