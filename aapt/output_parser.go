package aapt

import (
	"errors"
	"regexp"
	"strings"
)

func ParsePackageName(out string) (string, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		r := regexp.MustCompile(`package: name='([^']+)'`)
		matches := r.FindStringSubmatch(line)
		if matches != nil {
			return matches[1], nil
		}
	}

	return "", errors.New("package name not found")
}

func ParseTestRunner(out string) (string, error) {
	lines := strings.Split(out, "\n")
	startInstrumentationIdx := len(lines)
	for idx := 0; idx < len(lines); idx++ {
		line := lines[idx]
		regexInstrumentation := regexp.MustCompile(`E: instrumentation`)
		if regexInstrumentation.FindStringSubmatch(line) != nil {
			startInstrumentationIdx = idx
			break
		}
	}

	for idx := startInstrumentationIdx + 1; idx < len(lines); idx++ {
		line := lines[idx]
		regexName := regexp.MustCompile(`A: android:name\([^)]+\)="([^"]+)"`)
		matches := regexName.FindStringSubmatch(line)
		if matches != nil {
			return matches[1], nil
		}
	}

	return "", errors.New("Test runner not found in AndroidManifest")
}
