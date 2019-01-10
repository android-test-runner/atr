package adb

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type outputParser interface {
	ParseConnectedDeviceSerials(out string) []string
	ParseVersion(out string) (string, error)
	ParseDimensions(out string) (int, int, error)
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

func (outputParserImpl) ParseDimensions(out string) (int, int, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		r := regexp.MustCompile(`Physical size: ([0-9]+)x([0-9]+)`)
		matches := r.FindStringSubmatch(line)
		if matches != nil {
			// no conversion errors posible sinze we only match numbers
			width, _ := strconv.Atoi(matches[1])
			height, _ := strconv.Atoi(matches[2])
			return width, height, nil
		}
	}

	return 0, 0, errors.New("no dimesnions found")
}
