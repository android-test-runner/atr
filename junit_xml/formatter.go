package junit_xml

import (
	"encoding/xml"
	"fmt"
	"github.com/android-test-runner/atr/apks"
	"github.com/android-test-runner/atr/files"
	"github.com/android-test-runner/atr/result"
)

type Formatter interface {
	Format(result []result.Result, apk apks.Apk) (files.File, error)
}

type formatterImpl struct{}

func NewFormatter() Formatter {
	return formatterImpl{}
}

type skipped struct{}

type testcase struct {
	XMLName    xml.Name `xml:"testcase"`
	MethodName string   `xml:"name,attr"`
	ClassName  string   `xml:"classname,attr"`
	Failure    string   `xml:"failure,omitempty"`
	Error      string   `xml:"error,omitempty"`
	Skipped    *skipped `xml:"skipped,omitempty"`
	Time       string   `xml:"time,attr"`
}

type testsuite struct {
	XMLName     xml.Name   `xml:"testsuite"`
	Properties  string     `xml:"properties"`
	Name        string     `xml:"name,attr"`
	NumTests    int        `xml:"tests,attr"`
	NumFailures int        `xml:"failures,attr"`
	NumErrors   int        `xml:"errors,attr"`
	NumSkipped  int        `xml:"skipped,attr"`
	Time        string     `xml:"time,attr"`
	Timestamp   string     `xml:"timestamp,attr"`
	TestCases   []testcase `xml:"testcase"`
}

const (
	durationFormat = "%.3f"
)

func (formatterImpl) Format(results []result.Result, apk apks.Apk) (files.File, error) {
	var testCases []testcase
	numFailures := 0
	numErrors := 0
	numSkipped := 0
	totalTime := 0.0
	for _, r := range results {
		testCase := testcase{
			MethodName: r.Test.Method,
			ClassName:  r.Test.Class,
			Time:       fmt.Sprintf(durationFormat, r.Duration.Seconds()),
		}
		if r.Status == result.Errored {
			testCase.Error = r.Output
			numErrors += 1
		} else if r.Status == result.Failed {
			testCase.Failure = r.Output
			numFailures += 1
		} else if r.Status == result.Skipped {
			testCase.Skipped = &skipped{}
			numSkipped += 1
		}
		totalTime += r.Duration.Seconds()

		testCases = append(testCases, testCase)
	}

	testSuite := &testsuite{
		Name:        apk.PackageName,
		TestCases:   testCases,
		NumTests:    len(testCases),
		NumFailures: numFailures,
		NumErrors:   numErrors,
		NumSkipped:  numSkipped,
		Time:        fmt.Sprintf(durationFormat, totalTime),
		// Temporary fix to make test-store happy. It needs to be decided whether
		// to add this field here or to make test-store not depend on it.
		Timestamp:   "1970-01-01T00:00:00",
	}

	output, err := xml.MarshalIndent(testSuite, "", "    ")

	if err != nil {
		return files.File{}, err
	}

	f := files.File{
		Name:    "junit.xml",
		Content: string(output),
	}

	return f, nil
}
