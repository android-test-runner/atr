package output

import (
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
	"path/filepath"
)

type Writer interface {
	WriteFile(file files.File, device devices.Device) (string, error)
	WriteFileToRoot(file files.File) (string, error)
	RemoveDeviceDirectory(device devices.Device) error
	MakeDeviceDirectory(device devices.Device) (string, error)
	ToAbsolute(path string) string
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
	path := writer.ToAbsolute(file.EscapedName())
	return file.EscapedName(), writer.files.WriteFile(path, file.Content)
}

func (writer writerImpl) WriteFile(file files.File, device devices.Device) (string, error) {
	relativePath := writer.relativePath(device, file)
	return relativePath, writer.files.WriteFile(writer.ToAbsolute(relativePath), file.Content)
}

func (writer writerImpl) RemoveDeviceDirectory(device devices.Device) error {
	path := writer.ToAbsolute(writer.getDeviceDirectory(device))
	return writer.files.RemoveDirectory(path)
}

func (writer writerImpl) MakeDeviceDirectory(device devices.Device) (string, error) {
	deviceDirectory := writer.getDeviceDirectory(device)
	err := writer.files.MakeDirectory(writer.ToAbsolute(deviceDirectory))
	return deviceDirectory, err
}

func (writer writerImpl) relativePath(device devices.Device, file files.File) string {
	return filepath.Join(device.Serial, file.EscapedName())
}

func (writer writerImpl) ToAbsolute(path string) string {
	return filepath.Join(writer.rootDir, path)
}

func (writer writerImpl) getDeviceDirectory(device devices.Device) string {
	return device.Serial
}
