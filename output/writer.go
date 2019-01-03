package output

import (
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/files"
	"path/filepath"
)

type Writer interface {
	Write(files map[devices.Device][]files.File) error
	WriteFile(file files.File, device devices.Device) error
}

type writerImpl struct {
	rootDir string
	files   files.Files
}

func NewWriter(rootDirectory string) Writer {
	return writerImpl{
		rootDir: rootDirectory,
		files:   files.New(),
	}
}

func (writer writerImpl) Write(files map[devices.Device][]files.File) error {
	for device, filesForDevice := range files {
		err := writer.write(filesForDevice, device)
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer writerImpl) WriteFile(file files.File, device devices.Device) error {
	deviceDirectory, errDirectory := writer.ensureDeviceDirectory(device)
	if errDirectory != nil {
		return errDirectory
	}

	errFile := writer.files.WriteFile(deviceDirectory, file)
	if errFile != nil {
		return errFile
	}

	return nil
}

func (writer writerImpl) write(files []files.File, device devices.Device) error {
	deviceDirectory, err := writer.ensureDeviceDirectory(device)
	if err != nil {
		return err
	}
	for _, f := range files {
		err := writer.files.WriteFile(deviceDirectory, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer writerImpl) ensureDeviceDirectory(device devices.Device) (string, error) {
	deviceDirectory := filepath.Join(writer.rootDir, device.Serial)
	return deviceDirectory, writer.files.MakeDirectory(deviceDirectory)
}
