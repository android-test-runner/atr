package result

import (
	"github.com/ybonjour/atr/test"
	"time"
)

type Result struct {
	Test     test.Test
	Status   Status
	Output   string
	Duration time.Duration
	Extras   []Extra
}

func (result Result) equals(otherResult Result) bool {
	return result.Test == otherResult.Test &&
		result.Status == otherResult.Status &&
		result.Output == otherResult.Output &&
		result.Duration == otherResult.Duration &&
		result.areEqual(result.Extras, otherResult.Extras)
}

func (result Result) areEqual(slice1, slice2 []Extra) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

type Extra struct {
	Name  string
	Value string
	Type  ExtraType
}

type ExtraType string

const (
	File ExtraType = "file"
)

func (result Result) IsFailure() bool {
	return result.Status == Failed || result.Status == Errored
}

type Status int

const (
	Passed  Status = iota
	Failed  Status = iota
	Errored Status = iota
	Skipped Status = iota
)

func (status Status) toString() string {
	switch status {
	case Passed:
		return "Passed"
	case Failed:
		return "Failed"
	case Errored:
		return "Errored"
	case Skipped:
		return "Skipped"
	default:
		return "Unknown"
	}
}
