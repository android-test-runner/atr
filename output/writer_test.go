package output

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/files"
	"github.com/ybonjour/atr/mock_files"
	"testing"
)

func TestWrite(t *testing.T) {
	rootDir := "rootDir"
	file := files.File{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	files_mock := mock_files.NewMockFiles(ctrl)
	files_mock.EXPECT().MakeDirectory(rootDir)
	files_mock.EXPECT().WriteFile(rootDir, file).Return(nil)

	writer := writerImpl{
		rootDir: rootDir,
		files:   files_mock,
	}

	err := writer.Write([]files.File{file})

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
}
