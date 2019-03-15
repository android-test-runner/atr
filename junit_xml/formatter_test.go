package junit_xml

import (
	"fmt"
	"github.com/android-test-runner/atr/apks"
	"github.com/android-test-runner/atr/result"
	"github.com/android-test-runner/atr/test"
	"strings"
	"testing"
	"time"
)

func TestFormatsPassedTest(t *testing.T) {
	passedTest := result.Result{
		Test:     test.Test{Class: "TestClass", Method: "testMethod"},
		Status:   result.Passed,
		Duration: 1*time.Second + 500*time.Millisecond,
	}
	apk := apks.Apk{PackageName: "ch.yvu.atr"}

	xmlFile, _ := NewFormatter().Format([]result.Result{passedTest}, apk)

	expectedFilename := "junit.xml"
	if xmlFile.Name != expectedFilename {
		t.Error(fmt.Sprintf("Expected file name to be '%v' but it was '%v'", expectedFilename, xmlFile.Name))
	}
	expectedXml := `<testsuite name="ch.yvu.atr" tests="1" failures="0" errors="0" skipped="0" time="1.500"><properties></properties><testcase name="testMethod" classname="TestClass" time="1.500"></testcase></testsuite>`
	if removeWhitespaces(expectedXml) != removeWhitespaces(xmlFile.Content) {
		t.Error(fmt.Sprintf("Expected xml '%v' but got '%v'.", expectedXml, xmlFile.Content))
	}
}

func TestFormatsFailedTest(t *testing.T) {
	failedTest := result.Result{
		Test:     test.Test{Class: "TestClass", Method: "testMethod"},
		Status:   result.Failed,
		Output:   "failureOutput",
		Duration: 1*time.Second + 500*time.Millisecond,
	}
	apk := apks.Apk{PackageName: "ch.yvu.atr"}

	xmlFile, _ := NewFormatter().Format([]result.Result{failedTest}, apk)

	expectedXml := `<testsuite name="ch.yvu.atr" tests="1" failures="1" errors="0" skipped="0" time="1.500"><properties></properties><testcase name="testMethod" classname="TestClass" time="1.500"><failure>failureOutput</failure></testcase></testsuite>`
	if removeWhitespaces(expectedXml) != removeWhitespaces(xmlFile.Content) {
		t.Error(fmt.Sprintf("Expected xml '%v' but got '%v'.", expectedXml, xmlFile.Content))
	}
}

func TestFormatsErroredTest(t *testing.T) {
	failedTest := result.Result{
		Test:     test.Test{Class: "TestClass", Method: "testMethod"},
		Status:   result.Errored,
		Output:   "errorOutput",
		Duration: 1*time.Second + 500*time.Millisecond,
	}
	apk := apks.Apk{PackageName: "ch.yvu.atr"}

	xmlFile, _ := NewFormatter().Format([]result.Result{failedTest}, apk)

	expectedXml := `<testsuite name="ch.yvu.atr" tests="1" failures="0" errors="1" skipped="0" time="1.500"><properties></properties><testcase name="testMethod" classname="TestClass" time="1.500"><error>errorOutput</error></testcase></testsuite>`
	if removeWhitespaces(expectedXml) != removeWhitespaces(xmlFile.Content) {
		t.Error(fmt.Sprintf("Expected xml '%v' but got '%v'.", expectedXml, xmlFile.Content))
	}
}

func TestFormatsSkippedTest(t *testing.T) {
	skippedTest := result.Result{
		Test:     test.Test{Class: "TestClass", Method: "testMethod"},
		Status:   result.Skipped,
		Duration: 1*time.Second + 500*time.Millisecond,
	}
	apk := apks.Apk{PackageName: "ch.yvu.atr"}

	xmlFile, _ := NewFormatter().Format([]result.Result{skippedTest}, apk)

	expectedXml := `<testsuite name="ch.yvu.atr" tests="1" failures="0" errors="0" skipped="1" time="1.500"><properties></properties><testcase name="testMethod" classname="TestClass" time="1.500"><skipped></skipped></testcase></testsuite>`

	if removeWhitespaces(expectedXml) != removeWhitespaces(xmlFile.Content) {
		t.Error(fmt.Sprintf("Expected xml '%v' but got '%v'.", expectedXml, xmlFile.Content))
	}
}

func TestFormatsMultipleTests(t *testing.T) {
	test1 := result.Result{
		Test:     test.Test{Class: "TestClass1", Method: "testMethod1"},
		Status:   result.Passed,
		Duration: 1 * time.Second,
	}
	test2 := result.Result{
		Test:     test.Test{Class: "TestClass2", Method: "testMethod2"},
		Duration: 1 * time.Second,
	}
	apk := apks.Apk{PackageName: "ch.yvu.atr"}

	xmlFile, _ := NewFormatter().Format([]result.Result{test1, test2}, apk)

	expectedXml := `<testsuite name="ch.yvu.atr" tests="2" failures="0" errors="0" skipped="0" time="2.000"><properties></properties><testcase name="testMethod1" classname="TestClass1" time="1.000"></testcase><testcase name="testMethod2" classname="TestClass2" time="1.000"></testcase></testsuite>`
	if removeWhitespaces(expectedXml) != removeWhitespaces(xmlFile.Content) {
		t.Error(fmt.Sprintf("Expected xml '%v' but got '%v'.", expectedXml, xmlFile.Content))
	}
}

func removeWhitespaces(s string) string {
	output := s
	output = strings.Replace(output, " ", "", -1)
	output = strings.Replace(output, "\n", "", -1)
	return output
}
