package aapt

import (
	"fmt"
	"testing"
)

func TestParsesTestRunner(t *testing.T) {
	expectedTestRunner := "testRunner"
	outTemplate := "E: instrumentation (line=1)\n" +
		"A: android:name(0x01010003)=\"%v\" (Raw: \"%v\")"
	out := fmt.Sprintf(outTemplate, expectedTestRunner, expectedTestRunner)

	testRunner, err := newOutputParser().ParseTestRunner(out)

	verifyTestRunnner(expectedTestRunner, testRunner, err, t)
}

func TestParsesTestWithOtherLines(t *testing.T) {
	expectedTestRunner := "testRunner"
	outTemplate := "line before\n" +
		"E: instrumentation (line=1)\n" +
		"line between \n" +
		"A: android:name(0x01010003)=\"%v\" (Raw: \"%v\")\n" +
		"line after"
	out := fmt.Sprintf(outTemplate, expectedTestRunner, expectedTestRunner)

	testRunner, err := newOutputParser().ParseTestRunner(out)

	verifyTestRunnner(expectedTestRunner, testRunner, err, t)
}

func TestDoesNotParseTestRunnerWithNoInstrumentationElement(t *testing.T) {
	out := "A: android:name(0x01010003)=\"testRunner\" (Raw: \"testRunner\")"

	_, err := newOutputParser().ParseTestRunner(out)

	verifyTestRunnerNotFoundError(err, t)
}

func TestDoesNotParseTestRunnerWithNoNameAttribute(t *testing.T) {
	out := "E: instrumentation (line=1)"

	_, err := newOutputParser().ParseTestRunner(out)

	verifyTestRunnerNotFoundError(err, t)
}

func verifyTestRunnner(expected string, actual string, err error, t *testing.T) {
	if err != nil {
		t.Error(fmt.Sprintf("Got unexpected error: %v", err))
	}
	if expected != actual {
		t.Error(fmt.Sprintf("Got testrunner %v isntead of %v", actual, expected))
	}
}

func verifyTestRunnerNotFoundError(err error, t *testing.T) {
	if err == nil {
		t.Error("Expected 'test runner not found' error but didn't get one.")
	}
}
