package result

import (
	"bytes"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
	"html/template"
	"sort"
)

const cssTemplate = `
* {
	font-family: Arial;
}
pre {
	font-family: monospace;
	padding: 5px;
}
p.title {
	margin: 0px;
	padding: 5px;
	font-weight: bold;
}
video {
	width: 400px;
	padding: 5px;
}
ul.testResults {
	list-style-type:none;
	padding-left: 0;
}
ul.extras {
	padding-bottom: 5px;
}
.Passed {
	border-left: 5px solid green;
	border-bottom: 1px solid green;
}
.Failed {
	border-left: 5px solid red;
	border-bottom: 1px solid red;
}
.Errored {
	border-left: 5px solid red;
	border-bottom: 1px solid red;
}
.Skipped {
	border-left: 5px solid yellow;
	border-bottom: 1px solid yellow;
}
`

const htmlTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<title>ATR Results</title>
		<link href="{{ .ResultsCss }}" rel="stylesheet" />
	</head>
	<body>
		{{ range $testResult := .Results }}
			<h1>{{ $testResult.DeviceName }}</h1>
			<ul class="testResults">
			{{ range $result := $testResult.Results }}
				<li class="testResults {{ $result.Status }}">
				<p class="title">{{ $result.TestName }}: {{$result.Status}}</p>
				{{ if $result.Output }}
					<pre>{{ $result.Output }}</pre>
				{{ end }}
				{{ if $result.Video }}
					<video controls>
						<source src="{{$result.Video}}" type="video/mp4" />
						Your browser does not support the video tag.
					</video>
				{{ end }}
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
	Results    []resultsForDeviceHtml
	ResultsCss string
}

type resultsForDeviceHtml struct {
	DeviceName string
	Results    []resultHtml
}

type resultHtml struct {
	TestName  string
	Status    string
	IsFailure bool
	Output    string
	Video     string
	Extras    []extraHtml
}

type extraHtml struct {
	Name string
	Link string
}

type HtmlFormatter interface {
	FormatResults(map[devices.Device]TestResults) ([]files.File, error)
}

type htmlFormatterImpl struct{}

func NewHtmlFormatter() HtmlFormatter {
	return htmlFormatterImpl{}
}

func (formatter htmlFormatterImpl) FormatResults(resultsByDevice map[devices.Device]TestResults) ([]files.File, error) {
	parsedTemplate, templateError := template.New("html").Parse(htmlTemplate)
	if templateError != nil {
		return []files.File{}, templateError
	}

	cssFileName := "results.css"
	var content bytes.Buffer
	htmlOutput := formatter.toHtmlOutput(resultsByDevice, cssFileName)
	templateExecutionError := parsedTemplate.Execute(&content, htmlOutput)
	if templateExecutionError != nil {
		return []files.File{}, templateExecutionError
	}

	htmlFile := files.File{
		Name:    "results.html",
		Content: content.String(),
	}

	cssFile := files.File{
		Name:    cssFileName,
		Content: cssTemplate,
	}

	return []files.File{htmlFile, cssFile}, nil
}

func (formatter htmlFormatterImpl) toHtmlOutput(resultsByDevice map[devices.Device]TestResults, cssFileName string) outputHtml {
	resultsForDeviceHtmls := []resultsForDeviceHtml{}
	for device, testResults := range resultsByDevice {
		resultsHtml := []resultHtml{}
		for _, result := range testResults.Results {
			resultsHtml = append(resultsHtml, toHtmlResult(result))
		}
		sortResultHtmlsByStatus(resultsHtml)

		resultsAndDevice := resultsForDeviceHtml{
			DeviceName: device.Serial,
			Results:    resultsHtml,
		}

		resultsForDeviceHtmls = append(resultsForDeviceHtmls, resultsAndDevice)
	}

	return outputHtml{Results: resultsForDeviceHtmls, ResultsCss: cssFileName}
}

func toHtmlResult(result Result) resultHtml {
	htmlExtras := []extraHtml{}
	video := ""
	for _, extra := range result.Extras {
		if extra.Type == File {
			htmlExtras = append(htmlExtras, toHtmlExtra(extra))
		} else if extra.Type == Video {
			video = extra.Value
		}
	}

	output := ""
	if result.IsFailure() {
		output = result.Output
	}

	return resultHtml{
		TestName:  result.Test.FullName(),
		Status:    result.Status.toString(),
		IsFailure: result.IsFailure(),
		Output:    output,
		Video:     video,
		Extras:    htmlExtras,
	}
}

func toHtmlExtra(extra Extra) extraHtml {
	return extraHtml{
		Name: extra.Name,
		Link: extra.Value,
	}
}

type ByStatus []resultHtml

func (s ByStatus) Len() int           { return len(s) }
func (s ByStatus) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByStatus) Less(i, j int) bool { return toSortNumber(s[i]) < toSortNumber(s[j]) }

func sortResultHtmlsByStatus(resultHtmls []resultHtml) {
	sort.Sort(ByStatus(resultHtmls))
}

func toSortNumber(result resultHtml) int {
	// A failed test result shall be before a successful one
	if result.IsFailure {
		return 0
	} else {
		return 1
	}
}
