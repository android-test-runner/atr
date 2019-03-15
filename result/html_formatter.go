package result

import (
	"bytes"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
	"html/template"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<title>ATR Results</title>
		<style type="text/css">
			* {
				font-family: Arial;
			}
			pre {
				font-family: monospace;
			}
			p.title {
				margin: 0px;
				padding: 5px;
			}
			ul.testResults {
				list-style-type:none;
				padding-left: 0;
			}
			li.testResults {
				border: 1px solid black;
			}
			.Passed {
				background-color: green;
				color: white;
			}
			.Failed {
				background-color: red;
				color: white;
			}
			.Errored {
				background-color: red;
				color: white;
			}
			.Skipped {
				background-color: yellow;
			}
		</style>
	</head>
	<body>
		{{ range $testResult := .Results }}
			<h1>{{ $testResult.DeviceName }}</h1>
			<ul class="testResults">
			{{ range $result := $testResult.Results }}
				<li class="testResults">
				<p class="title {{ $result.Status }}">{{ $result.TestName }}: {{$result.Status}}</p>
				<pre>{{ $result.Output }}</pre>
				<ul class="extras">
					{{ range $extra := $result.Extras }}
						<li><a href="{{ $extra.Link }}">{{ $extra.Name }}</a></li>
					{{ end }}
				</ul>
				</li>
			{{ end }}
			</ul>
		{{ end }}
	</body>
</html>
`

type outputHtml struct {
	Results []resultsForDeviceHtml
}

type resultsForDeviceHtml struct {
	DeviceName string
	Results    []resultHtml
}

type resultHtml struct {
	TestName string
	Status   string
	Output   string
	Extras   []extraHtml
}

type extraHtml struct {
	Name string
	Link string
}

type HtmlFormatter interface {
	FormatResults(map[devices.Device]TestResults) (files.File, error)
}

type htmlFormatterImpl struct{}

func NewHtmlFormatter() HtmlFormatter {
	return htmlFormatterImpl{}
}

func (formatter htmlFormatterImpl) FormatResults(resultsByDevice map[devices.Device]TestResults) (files.File, error) {
	parsedTemplate, templateError := template.New("html").Parse(htmlTemplate)
	if templateError != nil {
		return files.File{}, templateError
	}

	var content bytes.Buffer
	htmlOutput := formatter.toHtmlOutput(resultsByDevice)
	templateExecutionError := parsedTemplate.Execute(&content, htmlOutput)
	if templateExecutionError != nil {
		return files.File{}, templateExecutionError
	}

	file := files.File{
		Name:    "results.html",
		Content: content.String(),
	}

	return file, nil
}

func (formatter htmlFormatterImpl) toHtmlOutput(resultsByDevice map[devices.Device]TestResults) outputHtml {
	resultsForDeviceHtmls := []resultsForDeviceHtml{}
	for device, testResults := range resultsByDevice {
		resultsHtml := []resultHtml{}
		for _, result := range testResults.Results {
			resultsHtml = append(resultsHtml, toHtmlResult(result))
		}

		resultsAndDevice := resultsForDeviceHtml{
			DeviceName: device.Serial,
			Results:    resultsHtml,
		}

		resultsForDeviceHtmls = append(resultsForDeviceHtmls, resultsAndDevice)
	}

	return outputHtml{Results: resultsForDeviceHtmls}
}

func toHtmlResult(result Result) resultHtml {
	htmlExtras := []extraHtml{}
	for _, extra := range result.Extras {
		if extra.Type == File {
			htmlExtras = append(htmlExtras, toHtmlExtra(extra))
		}
	}

	return resultHtml{
		TestName: result.Test.FullName(),
		Status:   result.Status.toString(),
		Output:   result.Output,
		Extras:   htmlExtras,
	}
}

func toHtmlExtra(extra Extra) extraHtml {
	return extraHtml{
		Name: extra.Name,
		Link: extra.Value,
	}
}
