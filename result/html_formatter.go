package result

import (
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
)

type HtmlFormatter interface {
	FormatResults(map[devices.Device]TestResults) (files.File, error)
}

type htmlFormatterImpl struct{}

func NewHtmlFormatter() HtmlFormatter {
	return htmlFormatterImpl{}
}

func (formatter htmlFormatterImpl) FormatResults(map[devices.Device]TestResults) (files.File, error) {
	file := files.File{
		Name:    "results.html",
		Content: "",
	}

	return file, nil
}
