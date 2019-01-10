package result

import (
	"fmt"
	"testing"
)

func TestGetErrorFromFailure(t *testing.T) {
	failedResult := Result{Status: Failed}
	testResults := TestResults{
		Results: []Result{failedResult},
	}

	errors := testResults.ErrorsFromFailures()

	if errors == nil {
		t.Error("Expected an error for a failed test but got none")
	}
}

func TestGetNoErrorFromPassedTest(t *testing.T) {
	failedResult := Result{Status: Passed}
	testResults := TestResults{
		Results: []Result{failedResult},
	}

	errors := testResults.ErrorsFromFailures()

	if errors != nil {
		t.Error(fmt.Sprintf("Expected no errors but got: '%v'", errors))
	}
}
