package files

import "strings"

type File struct {
	Name    string
	Content string
}

func (file File) EscapedName() string {
	return EscapeFileName(file.Name)
}

func EscapeFileName(name string) string {
	return strings.Replace(name, "#", "_", -1)
}
