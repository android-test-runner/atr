package result

import (
	"encoding/json"
	"github.com/ybonjour/atr/test"
)

type test_json struct {
	Class  string `json:"class"`
	Method string `json:"method"`
}

type extra_json struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type result_json struct {
	Test            test_json    `json:"test"`
	Status          string       `json:"status"`
	Output          string       `json:"output"`
	DurationSeconds float64      `json:"durationSeconds"`
	Extras          []extra_json `json:"extras"`
}

func ToJson(results []Result) (string, error) {
	results_json := []result_json{}

	for _, result := range results {
		results_json = append(results_json, toJsonResult(result))
	}

	b, err := json.Marshal(results_json)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func toJsonResult(result Result) result_json {
	extras_json := []extra_json{}
	for _, extra := range result.Extras {
		extras_json = append(extras_json, toJsonExtra(extra))
	}
	return result_json{
		Test:            toJsonTest(result.Test),
		Status:          result.Status.toString(),
		Output:          result.Output,
		DurationSeconds: result.Duration.Seconds(),
		Extras:          extras_json,
	}
}

func toJsonTest(test test.Test) test_json {
	return test_json{
		Class:  test.Class,
		Method: test.Method,
	}
}

func toJsonExtra(extra Extra) extra_json {
	return extra_json{
		Name:  extra.Name,
		Value: extra.Value,
		Type:  string(extra.Type),
	}
}
