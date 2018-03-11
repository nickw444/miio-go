package common

type EventNewDevice struct {
	Device Device
}

type EventNewMaskedDevice struct {
	DeviceID uint32
}

type EventExpiredDevice struct {
	Device Device
}

type EventUpdatePower struct {
	PowerState PowerState
}

type EventUpdateLight struct {
	Brightness int

	ColorMode int // 1: rgb mode, 2: color temperature mode, 3: hsv mode
	RGB       struct {
		Red   int
		Green int
		Blue  int
	}
	Hue        int
	Saturation int
}
