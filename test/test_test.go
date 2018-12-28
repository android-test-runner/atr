package test

import (
	"fmt"
	"testing"
)

func TestGetsFullName(t *testing.T) {
	test := Test{Class: "TestClass", Method: "testMethod"}

	fullName := test.FullName()

	expected := "TestClass#testMethod"
	if expected != fullName {
		t.Error(fmt.Sprintf("Fullname is %v instead of %v", fullName, expected))
	}
}
