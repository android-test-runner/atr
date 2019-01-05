package adb

import (
	"errors"
	"regexp"
	"strings"
)

type outputParser interface {
	ParseConnectedDeviceSerials(out string) []string
	ParseVersion(out string) (string, error)
}

type outputParserImpl struct{}

func newOutputParser() outputParser {
	return outputParserImpl{}
}

func (outputParserImpl) ParseVersion(out string) (string, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		r := regexp.MustCompile(`Android Debug Bridge version ([0-9.]+)`)
		matches := r.FindStringSubmatch(line)
		if matches != nil {
			return matches[1], nil
		}
	}
	return "", errors.New("version not found")
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
