package test_executor

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
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
	installer     Installer
	resultParser  result.Parser
	adb           adb.Adb
	testListeners []test_listener.TestListener
	jsonFormatter result.JsonFormatter
	writer        output.Writer
}

func NewExecutor(writer output.Writer, testListeners []test_listener.TestListener) Executor {
	return executorImpl{
		installer:     NewInstaller(),
		resultParser:  result.NewParser(),
		adb:           adb.New(),
		testListeners: testListeners,
		jsonFormatter: result.NewJsonFormatter(),
		writer:        writer,
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
	executor.beforeTestSuite(device)
	var results []result.Result
	for _, t := range testConfig.Tests {
		executor.beforeTest(t, device)
		testOutput, errTest, duration := executor.executeSingleTest(t, device, testConfig.TestApk.PackageName, testConfig.TestRunner)
		r := executor.resultParser.ParseFromOutput(t, errTest, testOutput, duration)
		extendedResult := executor.afterTest(r, device)
		results = append(results, extendedResult)

	}
	executor.afterTestSuite(device)

	return results
}

func (executor executorImpl) forAllTestListeners(f func(listener test_listener.TestListener)) {
	for _, listener := range executor.testListeners {
		f(listener)
	}
}

func (executor executorImpl) beforeTestSuite(device devices.Device) {
	executor.forAllTestListeners(func(listener test_listener.TestListener) {
		listener.BeforeTestSuite(device)
	})
}

func (executor executorImpl) afterTestSuite(device devices.Device) {
	executor.forAllTestListeners(func(listener test_listener.TestListener) {
		listener.AfterTestSuite(device)
	})
}

func (executor executorImpl) beforeTest(t test.Test, device devices.Device) {
	executor.forAllTestListeners(func(listener test_listener.TestListener) {
		listener.BeforeTest(t, device)
	})
}

func (executor executorImpl) afterTest(r result.Result, device devices.Device) result.Result {
	executor.forAllTestListeners(func(listener test_listener.TestListener) {
		extras := listener.AfterTest(r, device)
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
