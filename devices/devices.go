package devices

import (
	"errors"
	"fmt"
	"github.com/ybonjour/atr/adb"
	"regexp"
	"strconv"
	"strings"
)

type Device struct {
	Serial          string
	ScreenDimension ScreenDimension
}

type ScreenDimension struct {
	Width  int
	Height int
}

func (dimension ScreenDimension) ToString() string {
	return fmt.Sprintf("%vx%v", dimension.Width, dimension.Height)
}

func ParseToScreenDimension(input string) (ScreenDimension, error) {
	r := regexp.MustCompile(`^([0-9]+)x([0-9]+)$`)
	matches := r.FindStringSubmatch(input)
	if matches == nil {
		return ScreenDimension{}, errors.New(fmt.Sprintf("'%v' is not a screen dimension", input))
	}

	// no errors possible in conversion to int, because we only match numbers
	width, _ := strconv.Atoi(matches[1])
	height, _ := strconv.Atoi(matches[2])

	return ScreenDimension{Width: width, Height: height}, nil
}

type Devices interface {
	ConnectedDevices(includeDeviceSerials []string) ([]Device, error)
	ParseConnectedDevices(deviceDefinitions []string) ([]Device, error)
	AllConnectedDevices() ([]Device, error)
}

type devicesImpl struct {
	adb adb.Adb
}

func New() Devices {
	return devicesImpl{
		adb: adb.New(),
	}
}

func (d devicesImpl) AllConnectedDevices() ([]Device, error) {
	connectedDeviceSerials, err := d.adb.ConnectedDevices()
	if err != nil {
		return nil, err
	}
	return d.parseConnectedDevices(connectedDeviceSerials, connectedDeviceSerials), nil
}

func (d devicesImpl) ParseConnectedDevices(deviceDefinitions []string) ([]Device, error) {
	connectedDeviceSerials, err := d.adb.ConnectedDevices()
	if err != nil {
		return nil, err
	}
	return d.parseConnectedDevices(deviceDefinitions, connectedDeviceSerials), nil
}

func (d devicesImpl) parseConnectedDevices(deviceDefinitions []string, connectedDeviceSerials []string) []Device {
	parsedDevices := []Device{}
	isConnected := toMap(connectedDeviceSerials)

	for _, deviceDefinition := range deviceDefinitions {
		tokens := strings.Split(deviceDefinition, "@")
		var parsedDevice *Device
		if len(tokens) == 2 {
			parsedDevice = d.getDeviceWithManualScreenDimensions(tokens[0], tokens[1])
		} else if len(tokens) == 1 {
			parsedDevice = d.getDeviceWithDefaultScreenDimensions(tokens[0])
		}
		if parsedDevice != nil && isConnected[parsedDevice.Serial] {
			parsedDevices = append(parsedDevices, *parsedDevice)
		}
	}
	return parsedDevices
}

func (d devicesImpl) getDeviceWithManualScreenDimensions(serial string, screenDimensions string) *Device {
	parsedDimensions, err := ParseToScreenDimension(screenDimensions)
	if err != nil {
		return nil
	}
	return &Device{Serial: serial, ScreenDimension: parsedDimensions}
}

func (d devicesImpl) getDeviceWithDefaultScreenDimensions(serial string) *Device {
	width, height, err := d.adb.GetScreenDimensions(serial)
	if err != nil {
		return nil
	}
	return &Device{Serial: serial, ScreenDimension: ScreenDimension{Width: width, Height: height}}
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
