package apks

import "os"

type Files interface {
	CanAccess(path string) bool
}

type filesImpl struct{}

func NewFiles() Files {
	return filesImpl{}
}

func (filesImpl) CanAccess(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
