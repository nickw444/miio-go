package device

import (
	"github.com/nickw444/miio-go/simulator/capability"
)

type SimulatedPowerPlug struct {
	*BaseDevice
}

func NewSimulatedPowerPlug(baseDevice *BaseDevice) *SimulatedPowerPlug {
	baseDevice.AddCapability(&capability.Info{
		Model: "chuangmi.plug.m1",
	})
	baseDevice.AddCapability(&capability.Power{})
	return &SimulatedPowerPlug{
		BaseDevice: baseDevice,
	}
}
