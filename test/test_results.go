package test

import (
	"regexp"
	"strings"
)

type TestResult struct {
	Test      Test
	HasPassed bool
	Output    string
}

func isTestSuccessful(output string) bool {
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
