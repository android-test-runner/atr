package device

type Device struct {
	Serial string
}

func FromSerials(serials []string) []Device {
	devices := []Device{}
	for _, serial := range serials {
		d := Device{
			Serial: serial,
		}
		devices = append(devices, d)
	}
	return devices
}
