package result

import (
	"errors"
	"fmt"
	"github.com/ybonjour/atr/test"
	"testing"
)

func TestResultFromOutputPassed(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	okOutput := "OK (1 test)"

	result := NewResultParser().ParseFromOutput(testForResult, nil, okOutput)

	expected := Result{Test: testForResult, HasPassed: true, Output: okOutput}
	if expected != result {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputPassedWithMultilineOutput(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	okOutput := "Some Information\n...\nOK (1 test)"

	result := NewResultParser().ParseFromOutput(testForResult, nil, okOutput)

	expected := Result{Test: testForResult, HasPassed: true, Output: okOutput}
	if expected != result {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputWithError(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	err := errors.New("some error")
	okOutput := "OK (1 test)"

	result := NewResultParser().ParseFromOutput(testForResult, err, okOutput)

	expected := Result{Test: testForResult, HasPassed: false, Output: okOutput}
	if expected != result {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputWithFailureOutput(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	failureOutput := "Failure"

	result := NewResultParser().ParseFromOutput(testForResult, nil, failureOutput)

	expected := Result{Test: testForResult, HasPassed: false, Output: failureOutput}
	if expected != result {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}
