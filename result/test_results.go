package result

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/ybonjour/atr/devices"
)

type TestResults struct {
	Device     devices.Device
	Results    []Result
	SetupError error
}

func (tr TestResults) ErrorsFromFailures() error {
	var combinedErrors error
	for _, r := range tr.Results {
		if r.IsFailure() {
			combinedErrors = multierror.Append(combinedErrors, errors.New(fmt.Sprintf("Test '%v' failed on device '%v'", r.Test.FullName(), tr.Device)))
		}
	}

	return combinedErrors
}
