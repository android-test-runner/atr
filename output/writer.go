package output

import (
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/files"
	"path/filepath"
)

type Writer interface {
	GetDeviceDirectory(device devices.Device) (string, error)
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

func (writer writerImpl) WriteFileToRoot(file files.File) error {
	err := writer.files.MakeDirectory(writer.rootDir)
	if err != nil {
		return err
	}
	return writer.files.WriteFile(writer.rootDir, file)
}

func (writer writerImpl) WriteFile(file files.File, device devices.Device) error {
	deviceDirectory, errDirectory := writer.GetDeviceDirectory(device)
	if errDirectory != nil {
		return errDirectory
	}

	return writer.files.WriteFile(deviceDirectory, file)
}

func (writer writerImpl) GetDeviceDirectory(device devices.Device) (string, error) {
	deviceDirectory := filepath.Join(writer.rootDir, device.Serial)
	return deviceDirectory, writer.files.MakeDirectory(deviceDirectory)
}
