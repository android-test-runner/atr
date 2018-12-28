package test

import (
	"errors"
	"fmt"
	"testing"
)

func TestResultFromOutputPassed(t *testing.T) {
	test := Test{Class: "TestClass", Method: "testMethod"}
	okOutput := "OK (1 test)"

	result := ResultFromOutput(test, nil, okOutput)

	expected := Result{Test: test, HasPassed: true, Output: okOutput}
	if expected != result {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputPassedWithMultilineOutput(t *testing.T) {
	test := Test{Class: "TestClass", Method: "testMethod"}
	okOutput := "Some Information\n...\nOK (1 test)"

	result := ResultFromOutput(test, nil, okOutput)

	expected := Result{Test: test, HasPassed: true, Output: okOutput}
	if expected != result {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputWithError(t *testing.T) {
	test := Test{Class: "TestClass", Method: "testMethod"}
	err := errors.New("some error")
	okOutput := "OK (1 test)"

	result := ResultFromOutput(test, err, okOutput)

	expected := Result{Test: test, HasPassed: false, Output: okOutput}
	if expected != result {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputWithFailureOutput(t *testing.T) {
	test := Test{Class: "TestClass", Method: "testMethod"}
	failureOutput := "Failure"

	result := ResultFromOutput(test, nil, failureOutput)

	expected := Result{Test: test, HasPassed: false, Output: failureOutput}
	if expected != result {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}
