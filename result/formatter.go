package result

import (
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
)

type Formatter interface {
	FormatResults(map[devices.Device]TestResults) ([]files.File, error)
}
