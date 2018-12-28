package test

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/mock_files"
	"testing"
)

func TestParsesTests(t *testing.T) {
	unparsedTests := []string{"TestClass1#testMethod1", "TestClass2#testMethod2"}

	parsedTests := NewParser().Parse(unparsedTests)

	expected := []Test{
		{Class: "TestClass1", Method: "testMethod1"},
		{Class: "TestClass2", Method: "testMethod2"},
	}
	if !AreEqual(expected, parsedTests) {
		t.Error(fmt.Sprintf("Parsed Tests are %v instead of %v", parsedTests, expected))
	}
}

func TestParsesTestsFromFile(t *testing.T) {
	path := "tests.txt"
	unparsedTesets := []string{"TestClass1#testMethod1", "TestClass2#testMethod2"}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filesMock := mock_files.NewMockFiles(ctrl)
	filesMock.EXPECT().ReadLines(gomock.Eq(path)).Return(unparsedTesets, nil)
	parser := parserImpl{
		files: filesMock,
	}

	parsedTests, err := parser.ParseFromFile(path)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v", err))
	}
	expected := []Test{
		{Class: "TestClass1", Method: "testMethod1"},
		{Class: "TestClass2", Method: "testMethod2"},
	}
	if !AreEqual(expected, parsedTests) {
		t.Error(fmt.Sprintf("Parsed Tests are %v instead of %v", parsedTests, expected))
	}
}

func TestReturnsErrorIfFileCanNotBeRead(t *testing.T) {
	expectedError := errors.New("Error reading file")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filesMock := mock_files.NewMockFiles(ctrl)
	filesMock.EXPECT().ReadLines(gomock.Any()).Return(nil, expectedError)
	parser := parserImpl{
		files: filesMock,
	}

	_, err := parser.ParseFromFile("tests.txt")

	if expectedError != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedError, err))
	}
}

func AreEqual(slice1, slice2 []Test) bool {
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
