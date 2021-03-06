package files

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Files interface {
	CanAccess(path string) bool
	ReadLines(path string) ([]string, error)
	WriteFile(path string, content string) error
	MakeDirectory(directory string) error
	RemoveDirectory(directory string) error
}

type filesImpl struct{}

func New() Files {
	return filesImpl{}
}

func (filesImpl) CanAccess(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (files filesImpl) ReadLines(path string) ([]string, error) {
	if !files.CanAccess(path) {
		return nil, errors.New(fmt.Sprintf("Can not access file %v", path))
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func (files filesImpl) WriteFile(path string, content string) error {
	errDirectory := files.MakeDirectory(filepath.Dir(path))
	if errDirectory != nil {
		return errDirectory
	}

	return ioutil.WriteFile(path, []byte(content), 0644)
}

func (filesImpl) MakeDirectory(directory string) error {
	return os.MkdirAll(directory, os.ModePerm)
}

func (filesImpl) RemoveDirectory(directory string) error {
	return os.RemoveAll(directory)
}
