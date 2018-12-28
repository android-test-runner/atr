package test

import (
	"testing"
)

func TestParse(t *testing.T) {
	unparsedTests := []string{"TestClass#testMethod"}

	parsedTests := ParseTests(unparsedTests)

	expected := []Test{{Class: "TestClass", Method: "testMethod"}}
	if !AreEqual(expected, parsedTests) {
		t.Error("Unexpected result", expected, parsedTests)
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
