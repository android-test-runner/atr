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
	if !haveSameElements(expectedFiles, files) {
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
	if !haveSameElements(expectedFiles, files) {
		t.Error(fmt.Sprintf("Expected files '%v' but got '%v'", expectedFiles, files))
	}
}

func haveSameElements(slice1, slice2 []File) bool {
	return len(slice1) == len(slice2) && containsAll(slice1, slice2...)
}

func containsAll(haystack []File, needles ...File) bool {
	exists := map[File]bool{}
	for _, f := range haystack {
		exists[f] = true
	}
	for _, needle := range needles {
		if !exists[needle] {
			return false
		}
	}
	return true
}
