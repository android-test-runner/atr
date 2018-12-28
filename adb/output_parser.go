package adb

import (
	"regexp"
	"strings"
)

type OutputParser interface {
	ParseConnectedDeviceSerials(out string) []string
}

type outputParserImpl struct{}

func newOutputParser() OutputParser {
	return outputParserImpl{}
}

func (outputParserImpl) ParseConnectedDeviceSerials(out string) []string {
	serials := []string{}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		r := regexp.MustCompile(`^([^ ]+)\tdevice$`)
		matches := r.FindStringSubmatch(line)
		if matches != nil {
			serials = append(serials, matches[1])
		}
	}

	return serials
}
