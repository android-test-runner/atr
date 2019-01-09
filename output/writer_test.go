package output

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/files"
	"github.com/ybonjour/atr/mock_files"
	"testing"
)

func TestWriteFile(t *testing.T) {
	rootDir := "rootDir"
	file := files.File{}
	device := devices.Device{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filesMock := mock_files.NewMockFiles(ctrl)
	filesMock.EXPECT().MakeDirectory(rootDir)
	filesMock.EXPECT().WriteFile(rootDir, file).Return(nil)

	writer := writerImpl{
		rootDir: rootDir,
		files:   filesMock,
	}

	err := writer.WriteFile(file, device)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
}

func TestWriteFileToRoot(t *testing.T) {
	rootDir := "rootDir"
	file := files.File{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filesMock := mock_files.NewMockFiles(ctrl)
	filesMock.EXPECT().MakeDirectory(rootDir).Return(nil)
	filesMock.EXPECT().WriteFile(rootDir, file).Return(nil)

	writer := writerImpl{
		rootDir: rootDir,
		files:   filesMock,
	}

	err := writer.WriteFileToRoot(file)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
}
