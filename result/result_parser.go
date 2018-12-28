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
	return Result{
		Test:      test,
		HasPassed: err == nil && hasPassed(output),
		Output:    output,
	}
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
