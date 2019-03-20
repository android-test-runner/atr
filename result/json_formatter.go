package result

import (
	"encoding/json"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
	"github.com/android-test-runner/atr/test"
)

type testJson struct {
	Class  string `json:"class"`
	Method string `json:"method"`
}

type extraJson struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type resultJson struct {
	Test            testJson    `json:"test"`
	Status          string      `json:"status"`
	Output          string      `json:"output"`
	DurationSeconds float64     `json:"durationSeconds"`
	Extras          []extraJson `json:"extras"`
}

type testResultsJson struct {
	SetupError *string      `json:"setupError"`
	Results    []resultJson `json:"results"`
}

type JsonFormatter interface {
	FormatResults(map[devices.Device]TestResults) ([]files.File, error)
}

type jsonFormatterImpl struct{}

func NewJsonFormatter() JsonFormatter {
	return jsonFormatterImpl{}
}

func (parser jsonFormatterImpl) FormatResults(resultsByDevice map[devices.Device]TestResults) ([]files.File, error) {
	output, err := parser.parseResultsToString(resultsByDevice)
	file := files.File{
		Name:    "results.json",
		Content: output,
	}

	return []files.File{file}, err
}

func (jsonFormatterImpl) parseResultsToString(resultsByDevice map[devices.Device]TestResults) (string, error) {
	results_json := map[string]testResultsJson{}

	for device, results := range resultsByDevice {
		testResults := testResultsJson{
			Results: toJsonResults(results.Results),
		}
		if results.SetupError != nil {
			setupError := results.SetupError.Error()
			if setupError != "" {
				testResults.SetupError = &setupError
			}
		}
		results_json[device.Serial] = testResults
	}

	b, err := json.Marshal(results_json)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func toJsonResults(results []Result) []resultJson {
	resultJsons := []resultJson{}

	for _, result := range results {
		resultJsons = append(resultJsons, toJsonResult(result))
	}

	return resultJsons
}

func toJsonResult(result Result) resultJson {
	extras_json := []extraJson{}
	for _, extra := range result.Extras {
		extras_json = append(extras_json, toJsonExtra(extra))
	}
	return resultJson{
		Test:            toJsonTest(result.Test),
		Status:          result.Status.toString(),
		Output:          result.Output,
		DurationSeconds: result.Duration.Seconds(),
		Extras:          extras_json,
	}
}

func toJsonTest(test test.Test) testJson {
	return testJson{
		Class:  test.Class,
		Method: test.Method,
	}
}

func toJsonExtra(extra Extra) extraJson {
	return extraJson{
		Name:  extra.Name,
		Value: extra.Value,
		Type:  string(extra.Type),
	}
}
