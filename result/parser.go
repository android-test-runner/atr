package result

import (
	"github.com/android-test-runner/atr/test"
	"regexp"
	"strings"
	"time"
)

type Parser interface {
	ParseFromOutput(test test.Test, err error, output string, duration time.Duration) Result
}

type parserImpl struct{}

func NewParser() Parser {
	return parserImpl{}
}

func (parserImpl) ParseFromOutput(test test.Test, err error, output string, duration time.Duration) Result {
	status := getStatus(output, err)
	return Result{
		Test:     test,
		Status:   status,
		Output:   output,
		Duration: duration,
	}
}

func getStatus(output string, err error) Status {
	if err != nil {
		return Errored
	}
	// A test was successful if we find "OK (1 test)" in the output
	// A test was skipped if we find "OK (0 tests)" in the output
	// else it failed
	// This is needed because the am process does not fail if the test fails.
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		regexSkipped := regexp.MustCompile(`^\s*OK \(0 tests\)\s*$`)
		if regexSkipped.MatchString(line) {
			return Skipped
		}
		regexOk := regexp.MustCompile(`^\s*OK \(1 test\)\s*$`)
		if regexOk.MatchString(line) {
			return Passed
		}
	}

	return Failed
}
