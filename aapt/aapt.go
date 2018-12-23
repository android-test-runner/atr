package aapt

import (
	"errors"
	"github.com/ybonjour/atr/command"
	"os/exec"
	"regexp"
	"strings"
)

func PackageName(apkPath string) (string, error) {
	out, executeError := command.ExecuteOutput(exec.Command("aapt", "dump", "badging", apkPath))
	if executeError != nil {
		return "", executeError
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
