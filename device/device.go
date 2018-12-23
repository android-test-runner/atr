package device

import "github.com/ybonjour/atr/adb"

type Device struct {
	Serial string
}

func FromSerials(serials []string) []Device {
	var devices []Device
	for _, serial := range serials {
		d := Device{
			Serial: serial,
		}
		devices = append(devices, d)
	}
	return devices
}

func AllConnectedDevices() ([]Device, error) {
	deviceSerials, err := adb.ConnectedDevices()
	if err != nil {
		return nil, err
	}

	return FromSerials(deviceSerials), nil
}
