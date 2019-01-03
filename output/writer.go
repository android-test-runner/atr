package output

import "github.com/ybonjour/atr/files"

type Writer interface {
	Write(files []files.File) error
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

func (writer writerImpl) Write(files []files.File) error {
	err := writer.files.MakeDirectory(writer.rootDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		err := writer.files.WriteFile(writer.rootDir, f)
		if err != nil {
			return err
		}
	}
	return nil
}
