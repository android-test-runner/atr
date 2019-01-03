package files

import (
	"fmt"
	"testing"
)

func TestConvertsLabeledContentToFiles(t *testing.T) {
	label := "label"
	content := "content"
	labeledContent := map[string]string{label: content}
	fileType := "fileType"

	files := ToFiles(labeledContent, fileType)

	expectedFiles := []File{{Name: label, Type: fileType, Content: content}}
	if !AreEqual(expectedFiles, files) {
		t.Error(fmt.Sprintf("Expected files '%v' but got '%v'", expectedFiles, files))
	}
}

func TestConvertsMultipleLabeledContentToFiles(t *testing.T) {
	label1 := "label1"
	label2 := "label2"
	content1 := "content1"
	content2 := "content2"
	labeledContent := map[string]string{label1: content1, label2: content2}
	fileType := "fileType"

	files := ToFiles(labeledContent, fileType)

	expectedFiles := []File{{Name: label1, Type: fileType, Content: content1}, {Name: label2, Type: fileType, Content: content2}}
	if !AreEqual(expectedFiles, files) {
		t.Error(fmt.Sprintf("Expected files '%v' but got '%v'", expectedFiles, files))
	}
}

func AreEqual(slice1, slice2 []File) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}
