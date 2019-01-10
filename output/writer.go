package output

import (
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/files"
	"path/filepath"
)

type Writer interface {
	WriteFile(file files.File, device devices.Device) (string, error)
	WriteFileToRoot(file files.File) (string, error)
	RemoveDeviceDirectory(device devices.Device) error
	MakeDeviceDirectory(device devices.Device) (string, error)
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

func (writer writerImpl) WriteFileToRoot(file files.File) (string, error) {
	path := writer.toAbsolute(file.Name)
	return file.Name, writer.files.WriteFile(path, file.Content)
}

func (writer writerImpl) WriteFile(file files.File, device devices.Device) (string, error) {
	relativePath := writer.relativePath(device, file)
	return relativePath, writer.files.WriteFile(writer.toAbsolute(relativePath), file.Content)
}

func (writer writerImpl) RemoveDeviceDirectory(device devices.Device) error {
	path := writer.toAbsolute(writer.getDeviceDirectory(device))
	return writer.files.RemoveDirectory(path)
}

func (writer writerImpl) MakeDeviceDirectory(device devices.Device) (string, error) {
	deviceDirectory := writer.toAbsolute(writer.getDeviceDirectory(device))
	err := writer.files.MakeDirectory(deviceDirectory)
	return deviceDirectory, err
}

func (writer writerImpl) relativePath(device devices.Device, file files.File) string {
	return filepath.Join(device.Serial, file.Name)
}

func (writer writerImpl) toAbsolute(path string) string {
	return filepath.Join(writer.rootDir, path)
}

func (writer writerImpl) getDeviceDirectory(device devices.Device) string {
	return device.Serial
}
