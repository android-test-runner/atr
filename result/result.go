package result

import "github.com/ybonjour/atr/test"

type Result struct {
	Test       test.Test
	HasPassed  bool
	WasSkipped bool
	Status     Status
	Output     string
}

type Status int

const (
	Passed  Status = iota
	Failed  Status = iota
	Skipped Status = iota
)
