package result

import (
	"errors"
	"fmt"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/test"
	"strings"
	"testing"
	"time"
)

func TestConvertResultsToJson(t *testing.T) {
	result1 := Result{
		Test:     test.Test{Class: "TestClass", Method: "testMethod"},
		Status:   Failed,
		Output:   "Failure output",
		Duration: 2*time.Second + 300*time.Millisecond,
		Extras: []Extra{
			{Name: "extrafile1", Value: "/path/to/extrafile1", Type: File},
			{Name: "extrafile2", Value: "/path/to/extrafile2", Type: File},
		},
	}
	result2 := Result{
		Test:     test.Test{Class: "TestClass", Method: "testMethod1"},
		Status:   Passed,
		Duration: 2 * time.Second,
	}
	device := devices.Device{Serial: "deviceSerial"}
	testResults := TestResults{
		Results: []Result{result1, result2},
	}

	jsonFile, err := NewJsonFormatter().FormatResults(map[devices.Device]TestResults{device: testResults})

	if err != nil {
		t.Error(fmt.Printf("Did not expect an error but got '%v'", err))
	}
	expectedJson := `{
		"deviceSerial": {
			"setupError": null,
			"results": [{
				"test": {"class": "TestClass", "method": "testMethod"},
				"status": "Failed",
				"output": "Failure output",
				"durationSeconds": 2.3, 
				"extras": [
					{ "name": "extrafile1", "value": "/path/to/extrafile1", "type": "file"}, 
					{ "name": "extrafile2", "value": "/path/to/extrafile2", "type": "file"}
				] 
			},
			{
				"test": {"class": "TestClass", "method": "testMethod1"},
				"status": "Passed",
				"output": "",
				"durationSeconds": 2,
				"extras": []
			}]
		}
	}`
	if removeWhitespaces(expectedJson) != removeWhitespaces(jsonFile.Content) {
		t.Error(fmt.Printf("Expected json '%v' but got '%v'", removeWhitespaces(expectedJson), removeWhitespaces(jsonFile.Content)))
	}
}

func TestCanConvertSetupError(t *testing.T) {
	setupError := errors.New("setup error")
	testResults := TestResults{SetupError: setupError}
	device := devices.Device{Serial: "deviceSerial"}

	jsonFile, err := NewJsonFormatter().FormatResults(map[devices.Device]TestResults{device: testResults})
	if err != nil {
		t.Error(fmt.Printf("Did not expect an error but got '%v", err))
	}
	expectedJson := `{
		"deviceSerial": {
			"setupError": "setup error",
			"results": []
		}
	}`
	if removeWhitespaces(expectedJson) != removeWhitespaces(jsonFile.Content) {
		t.Error(fmt.Printf("Expected json '%v' but got '%v'", removeWhitespaces(expectedJson), removeWhitespaces(jsonFile.Content)))
	}
}

func removeWhitespaces(s string) string {
	output := s
	output = strings.Replace(output, " ", "", -1)
	output = strings.Replace(output, "\n", "", -1)
	output = strings.Replace(output, "\t", "", -1)
	return output
}
