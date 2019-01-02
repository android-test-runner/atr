package junit_xml

import (
	"encoding/xml"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/result"
)

type Formatter interface {
	Format(result []result.Result, apk apks.Apk) (string, error)
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
}

type testsuite struct {
	XMLName     xml.Name   `xml:"testsuite"`
	Properties  string     `xml:"properties"`
	Name        string     `xml:"name,attr"`
	NumTests    int        `xml:"tests,attr"`
	NumFailures int        `xml:"failures,attr"`
	NumErrors   int        `xml:"errors,attr"`
	NumSkipped  int        `xml:"skipped,attr"`
	TestCases   []testcase `xml:"testcase"`
}

func (formatterImpl) Format(results []result.Result, apk apks.Apk) (string, error) {
	var testCases []testcase
	numFailures := 0
	numErrors := 0
	numSkipped := 0
	for _, r := range results {
		testCase := testcase{
			MethodName: r.Test.Method,
			ClassName:  r.Test.Class,
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

		testCases = append(testCases, testCase)
	}

	testSuite := &testsuite{
		Name:        apk.PackageName,
		TestCases:   testCases,
		NumTests:    len(testCases),
		NumFailures: numFailures,
		NumErrors:   numErrors,
		NumSkipped:  numSkipped,
	}

	output, err := xml.MarshalIndent(testSuite, "", "    ")

	if err != nil {
		return "", err
	}

	return string(output), nil
}
