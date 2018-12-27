package aapt

import (
	"errors"
	"github.com/ybonjour/atr/command"
	"os/exec"
	"regexp"
	"strings"
)

func PackageName(apkPath string) (string, error) {
	out, err := command.ExecuteOutput(exec.Command("aapt", "dump", "badging", apkPath))
	if err != nil {
		return "", err
	}

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

func TestRunner(apkPath string) (string, error) {
	arguments := []string{
		"dunmp",
		"xmltree",
		apkPath,
		"AndroidManifest.xml",
	}

	out, err := command.ExecuteOutput(exec.Command("aapt", arguments...))
	if err != nil {
		return "", err
	}

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

	return "", errors.New("No test runner found in AndroidManifest")
}
