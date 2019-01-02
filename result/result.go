package result

import "github.com/ybonjour/atr/test"

type Result struct {
	Test       test.Test
	HasPassed  bool
	WasSkipped bool
	Output     string
}
