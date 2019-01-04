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
}

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
