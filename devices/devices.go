package devices

import (
	"github.com/ybonjour/atr/adb"
)

type Device struct {
	Serial string
}

type Devices interface {
	ConnectedDevices(includeDeviceSerials []string) ([]Device, error)
}

type devicesImpl struct {
	adb adb.Adb
}

func New() Devices {
	return devicesImpl{
		adb: adb.New(),
	}
}

func (d devicesImpl) ConnectedDevices(includeDeviceSerials []string) ([]Device, error) {
	allDevices, err := d.allConnectedDevices()
	if err != nil {
		return nil, err
	}

	if includeDeviceSerials == nil {
		return allDevices, nil
	}

	includedSerials := toMap(includeDeviceSerials)

	var filteredDevices []Device
	for _, device := range allDevices {
		if includedSerials[device.Serial] {
			filteredDevices = append(filteredDevices, device)
		}
	}

	return filteredDevices, nil
}

func (d devicesImpl) fromSerials(serials []string) []Device {
	var devices []Device
	for _, serial := range serials {
		d := Device{
			Serial: serial,
		}
		devices = append(devices, d)
	}
	return devices
}

func (d devicesImpl) allConnectedDevices() ([]Device, error) {
	deviceSerials, err := d.adb.ConnectedDevices()
	if err != nil {
		return nil, err
	}

	return d.fromSerials(deviceSerials), nil
}

func toMap(keys []string) map[string]bool {
	result := map[string]bool{}
	for _, key := range keys {
		result[key] = true
	}
	return result
}
