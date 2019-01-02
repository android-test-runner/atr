package output

type File struct {
	Name    string
	Content string
	Type    string
}

func ToFiles(labeledContent map[string]string, fileType string) []File {
	files := []File{}
	for label, content := range labeledContent {
		files = append(files, File{Name: label, Type: fileType, Content: content})
	}
	return files
}
