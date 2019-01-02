package junit_xml

import (
	"fmt"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"strings"
	"testing"
)

func TestFormatsPassedTest(t *testing.T) {
	passedTest := result.Result{
		Test:   test.Test{Class: "TestClass", Method: "testMethod"},
		Status: result.Passed,
	}
	apk := apks.Apk{PackageName: "ch.yvu.atr"}

	xmlOutput, _ := NewFormatter().Format([]result.Result{passedTest}, apk)

	expectedXml := `<testsuite name="ch.yvu.atr" tests="1" failures="0"><properties></properties><testcase name="testMethod" classname="TestClass"></testcase></testsuite>`
	if removeWhitespaces(expectedXml) != removeWhitespaces(xmlOutput) {
		t.Error(fmt.Sprintf("Expected xml '%v' but got '%v'.", expectedXml, xmlOutput))
	}
}

func TestFormatsFailedTest(t *testing.T) {
	passedTest := result.Result{
		Test:   test.Test{Class: "TestClass", Method: "testMethod"},
		Status: result.Failed,
		Output: "failureOutput",
	}
	apk := apks.Apk{PackageName: "ch.yvu.atr"}

	xmlOutput, _ := NewFormatter().Format([]result.Result{passedTest}, apk)

	expectedXml := `<testsuite name="ch.yvu.atr" tests="1" failures="1"><properties></properties><testcase name="testMethod" classname="TestClass"><failure>failureOutput</failure></testcase></testsuite>`
	if removeWhitespaces(expectedXml) != removeWhitespaces(xmlOutput) {
		t.Error(fmt.Sprintf("Expected xml '%v' but got '%v'.", expectedXml, xmlOutput))
	}
}

func TestFormatsMultipleTests(t *testing.T) {
	test1 := result.Result{
		Test:   test.Test{Class: "TestClass1", Method: "testMethod1"},
		Status: result.Passed,
	}
	test2 := result.Result{
		Test: test.Test{Class: "TestClass2", Method: "testMethod2"},
	}
	apk := apks.Apk{PackageName: "ch.yvu.atr"}

	xmlOutput, _ := NewFormatter().Format([]result.Result{test1, test2}, apk)

	expectedXml := `<testsuite name="ch.yvu.atr" tests="2" failures="0"><properties></properties><testcase name="testMethod1" classname="TestClass1"></testcase><testcase name="testMethod2" classname="TestClass2"></testcase></testsuite>`
	if removeWhitespaces(expectedXml) != removeWhitespaces(xmlOutput) {
		t.Error(fmt.Sprintf("Expected xml '%v' but got '%v'.", expectedXml, xmlOutput))
	}
}

func removeWhitespaces(s string) string {
	output := s
	output = strings.Replace(output, " ", "", -1)
	output = strings.Replace(output, "\n", "", -1)
	return output
}
