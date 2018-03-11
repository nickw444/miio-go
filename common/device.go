package common

import "github.com/nickw444/miio-go/subscription"

type DeviceInfo struct {
	FirmwareVersion string `json:"fw_ver"`
	HardwareVersion string `json:"hw_ver"`
	MacAddress      string `json:"mac"`
	Model           string `json:"model"`
}

type Device interface {
	subscription.SubscriptionTarget

	ID() uint32
	GetLabel() (string, error)
	GetInfo() (DeviceInfo, error)
	GetToken() []byte
}
