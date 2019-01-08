package result

import (
	"errors"
	"fmt"
	"github.com/ybonjour/atr/test"
	"testing"
	"time"
)

func TestResultFromOutputPassed(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	okOutput := "OK (1 test)"
	var duration time.Duration = 42

	result := NewParser().ParseFromOutput(testForResult, nil, okOutput, duration)

	expected := Result{Test: testForResult, Status: Passed, Output: okOutput, Duration: duration}
	if !expected.equals(result) {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputSkipped(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	skippedOutput := "OK (0 tests)"
	var duration time.Duration = 42

	result := NewParser().ParseFromOutput(testForResult, nil, skippedOutput, duration)

	expected := Result{Test: testForResult, Status: Skipped, Output: skippedOutput, Duration: duration}
	if !expected.equals(result) {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputPassedWithMultilineOutput(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	okOutput := "Some Information\n...\nOK (1 test)"
	var duration time.Duration = 42

	result := NewParser().ParseFromOutput(testForResult, nil, okOutput, duration)

	expected := Result{Test: testForResult, Status: Passed, Output: okOutput, Duration: duration}
	if !expected.equals(result) {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputErrored(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	err := errors.New("some error")
	okOutput := "OK (1 test)"
	var duration time.Duration = 42

	result := NewParser().ParseFromOutput(testForResult, err, okOutput, duration)

	expected := Result{Test: testForResult, Status: Errored, Output: okOutput, Duration: duration}
	if !expected.equals(result) {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}

func TestResultFromOutputFailed(t *testing.T) {
	testForResult := test.Test{Class: "TestClass", Method: "testMethod"}
	failureOutput := "Failure"

	result := NewParser().ParseFromOutput(testForResult, nil, failureOutput, 0)

	expected := Result{Test: testForResult, Status: Failed, Output: failureOutput}
	if !expected.equals(result) {
		t.Error(fmt.Sprintf("Test result is %v instead of %v", result, expected))
	}
}
