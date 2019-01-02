package result

import (
	"github.com/ybonjour/atr/test"
	"regexp"
	"strings"
)

type ResultParser interface {
	ParseFromOutput(test test.Test, err error, output string) Result
}

type resultParserImpl struct{}

func NewResultParser() ResultParser {
	return resultParserImpl{}
}

func (resultParserImpl) ParseFromOutput(test test.Test, err error, output string) Result {
	wasSkipped := wasSkipped(output)
	hasPassed := hasPassed(output)
	return Result{
		Test:       test,
		WasSkipped: wasSkipped,
		HasPassed:  wasSkipped || (err == nil && hasPassed),
		Output:     output,
	}
}

func wasSkipped(output string) bool {
	// A test was successful if we find "OK (0 tests)" in the output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		regexSkipped := regexp.MustCompile(`^OK \(0 tests\)$`)
		if regexSkipped.MatchString(line) {
			return true
		}
	}

	return false
}

func hasPassed(output string) bool {
	// A test was successful if we find "OK (1 test)" in the output
	// This is needed because the am process does not fail if the test fails.
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		regexOk := regexp.MustCompile(`^OK \(1 test\)$`)
		if regexOk.MatchString(line) {
			return true
		}
	}

	return false
}
