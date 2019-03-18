package test_executor

import (
	"fmt"
	"github.com/android-test-runner/atr/adb"
	"github.com/android-test-runner/atr/apks"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/output"
	"github.com/android-test-runner/atr/result"
	"github.com/android-test-runner/atr/test"
	"github.com/android-test-runner/atr/test_listener"
	"github.com/hashicorp/go-multierror"
	"sync"
	"time"
)

type Config struct {
	Apk               apks.Apk
	TestApk           apks.Apk
	TestRunner        string
	Tests             []test.Test
	OutputFolder      string
	DisableAnimations bool
}

type Executor interface {
	Execute(config Config, devices []devices.Device) error
}

type executorImpl struct {
	installer            Installer
	resultParser         result.Parser
	adb                  adb.Adb
	testListenersFactory test_listener.Factory
	jsonFormatter        result.JsonFormatter
	htmlFormatter        result.HtmlFormatter
	writer               output.Writer
}

func NewExecutor(writer output.Writer, testListenersFactory test_listener.Factory) Executor {
	return executorImpl{
		installer:            NewInstaller(),
		resultParser:         result.NewParser(),
		adb:                  adb.New(),
		testListenersFactory: testListenersFactory,
		jsonFormatter:        result.NewJsonFormatter(),
		htmlFormatter:        result.NewHtmlFormatter(),
		writer:               writer,
	}
}

func (executor executorImpl) Execute(config Config, targetDevices []devices.Device) error {
	resultsChannel := make(chan result.TestResults, len(targetDevices))

	var wg sync.WaitGroup
	wg.Add(len(targetDevices))
	for _, targetDevice := range targetDevices {
		go func(d devices.Device) {
			results, err := executor.executeOnDevice(config, d)
			resultsChannel <- result.TestResults{Device: d, Results: results, SetupError: err}

			wg.Done()
		}(targetDevice)
	}
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	var allErrors error
	resultsByDevice := map[devices.Device]result.TestResults{}

	for results := range resultsChannel {
		if results.SetupError != nil {
			allErrors = multierror.Append(results.SetupError)
		}

		errorsFromFailure := results.ErrorsFromFailures()
		if errorsFromFailure != nil {
			allErrors = multierror.Append(allErrors, errorsFromFailure)
		}

		resultsByDevice[results.Device] = results
	}

	errJson := executor.storeResultsAsJson(resultsByDevice)
	if errJson != nil {
		fmt.Printf("Could not write results.json: '%v'", errJson)
	}

	errHtml := executor.storeResultsAsHtml(resultsByDevice)
	if errHtml != nil {
		fmt.Printf("Could not write results.html: '%v'", errHtml)
	}

	return allErrors
}

func (executor executorImpl) storeResultsAsJson(resultsByDevice map[devices.Device]result.TestResults) error {
	file, errFormat := executor.jsonFormatter.FormatResults(resultsByDevice)
	if errFormat != nil {
		return errFormat
	}

	_, errWrite := executor.writer.WriteFileToRoot(file)

	return errWrite
}

func (executor executorImpl) storeResultsAsHtml(resultsByDevice map[devices.Device]result.TestResults) error {
	files, errFormat := executor.htmlFormatter.FormatResults(resultsByDevice)
	if errFormat != nil {
		return errFormat
	}

	for _, file := range files {
		_, errWrite := executor.writer.WriteFileToRoot(file)
		if errWrite != nil {
			return errWrite
		}
	}

	return nil
}

func (executor executorImpl) executeOnDevice(config Config, device devices.Device) ([]result.Result, error) {
	installError := executor.reinstallApks(config, device)
	if installError != nil {
		return nil, installError
	}

	removeError := executor.writer.RemoveDeviceDirectory(device)
	if removeError != nil {
		return nil, removeError
	}
	if config.DisableAnimations {
		disableAnimationsError := executor.adb.DisableAnimations(device.Serial)
		if disableAnimationsError != nil {
			return nil, disableAnimationsError
		}
	}

	return executor.executeTests(config, device), nil
}

func (executor executorImpl) reinstallApks(config Config, device devices.Device) error {
	apkInstallError := executor.installer.Reinstall(config.Apk, device)
	if apkInstallError != nil {
		return apkInstallError
	}
	testApkInstallError := executor.installer.Reinstall(config.TestApk, device)
	if testApkInstallError != nil {
		return testApkInstallError
	}

	return nil
}

func (executor executorImpl) executeTests(testConfig Config, device devices.Device) []result.Result {
	testListeners := executor.testListenersFactory.ForDevice(device)
	executor.beforeTestSuite(testListeners)
	var results []result.Result
	for _, t := range testConfig.Tests {
		executor.beforeTest(testListeners, t)
		testOutput, errTest, duration := executor.executeSingleTest(t, device, testConfig.TestApk.PackageName, testConfig.TestRunner)
		r := executor.resultParser.ParseFromOutput(t, errTest, testOutput, duration)
		extendedResult := executor.afterTest(testListeners, r)
		results = append(results, extendedResult)

	}
	executor.afterTestSuite(testListeners)

	return results
}

func (executor executorImpl) beforeTestSuite(testListeners []test_listener.TestListener) {
	executor.forAll(testListeners, func(listener test_listener.TestListener) {
		listener.BeforeTestSuite()
	})
}

func (executor executorImpl) afterTestSuite(testListeners []test_listener.TestListener) {
	executor.forAll(testListeners, func(listener test_listener.TestListener) {
		listener.AfterTestSuite()
	})
}

func (executor executorImpl) beforeTest(testListeners []test_listener.TestListener, t test.Test) {
	executor.forAll(testListeners, func(listener test_listener.TestListener) {
		listener.BeforeTest(t)
	})
}

func (executor executorImpl) afterTest(testListeners []test_listener.TestListener, r result.Result) result.Result {
	executor.forAll(testListeners, func(listener test_listener.TestListener) {
		extras := listener.AfterTest(r)
		r.Extras = append(r.Extras, extras...)
	})

	return r
}

func (executor executorImpl) executeSingleTest(t test.Test, device devices.Device, testPackage string, testRunner string) (string, error, time.Duration) {
	start := time.Now()
	testOutput, err := executor.adb.ExecuteTest(testPackage, testRunner, t.FullName(), device.Serial)
	duration := time.Since(start)
	return testOutput, err, duration
}

func (executor executorImpl) forAll(testListeners []test_listener.TestListener, f func(listener test_listener.TestListener)) {
	for _, listener := range testListeners {
		f(listener)
	}
}
