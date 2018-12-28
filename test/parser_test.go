package test

import (
	"fmt"
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

func TestGetsFullName(t *testing.T) {
	test := Test{Class: "TestClass", Method: "testMethod"}

	fullName := test.FullName()

	expected := "TestClass#testMethod"
	if expected != fullName {
		t.Error(fmt.Sprintf("Fullname is %v instead of %v", fullName, expected))
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
