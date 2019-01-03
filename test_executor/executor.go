package test_executor

import (
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/files"
	"github.com/ybonjour/atr/logcat"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/screen_recorder"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
	"sync"
	"time"
)

type Config struct {
	Apk          apks.Apk
	TestApk      apks.Apk
	TestRunner   string
	Tests        []test.Test
	OutputFolder string
}

type Executor interface {
	Execute(config Config, devices []devices.Device) (map[devices.Device][]result.Result, error)
}

type executorImpl struct {
	installer     Installer
	resultParser  result.Parser
	adb           adb.Adb
	testListeners []test_listener.TestListener
	writer        output.Writer
	files         files.Files
}

func NewExecutor(writer output.Writer) Executor {
	return executorImpl{
		installer:     NewInstaller(),
		resultParser:  result.NewParser(),
		adb:           adb.New(),
		testListeners: []test_listener.TestListener{logcat.NewLogcatListener(writer), screen_recorder.NewScreenRecorderListener(writer)},
		writer:        writer,
		files:         files.New(),
	}
}

type testResults struct {
	Results []result.Result
	Error   error
	Device  devices.Device
}

func (executor executorImpl) Execute(config Config, targetDevices []devices.Device) (map[devices.Device][]result.Result, error) {
	resultsChannel := make(chan testResults, len(targetDevices))

	var wg sync.WaitGroup
	wg.Add(len(targetDevices))
	for _, targetDevice := range targetDevices {
		go func(d devices.Device) {
			results, err := executor.executeOnDevice(config, d)
			resultsChannel <- testResults{Results: results, Error: err, Device: d}
			wg.Done()
		}(targetDevice)
	}
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	resultsByDevice := map[devices.Device][]result.Result{}
	for r := range resultsChannel {
		if r.Error != nil {
			return nil, r.Error
		}
		resultsByDevice[r.Device] = r.Results
	}

	return resultsByDevice, nil
}

func (executor executorImpl) executeOnDevice(config Config, device devices.Device) ([]result.Result, error) {
	installError := executor.reinstallApks(config, device)
	if installError != nil {
		return nil, installError
	}
	directory, directoryError := executor.writer.GetDeviceDirectory(device)
	if directoryError != nil {
		return nil, directoryError
	}
	removeError := executor.files.RemoveDirectory(directory)
	if removeError != nil {
		return nil, removeError
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
	var results []result.Result
	executor.beforeTestSuite(device)
	for _, t := range testConfig.Tests {
		executor.beforeTest(t)
		testOutput, errTest, duration := executor.executeSingleTest(t, device, testConfig.TestApk.PackageName, testConfig.TestRunner)
		r := executor.resultParser.ParseFromOutput(t, errTest, testOutput, duration)
		executor.afterTest(r)

		results = append(results, r)
	}
	return results
}

func (executor executorImpl) beforeTestSuite(device devices.Device) {
	for _, listener := range executor.testListeners {
		listener.BeforeTestSuite(device)
	}
}

func (executor executorImpl) beforeTest(t test.Test) {
	for _, listener := range executor.testListeners {
		listener.BeforeTest(t)
	}
}

func (executor executorImpl) afterTest(r result.Result) {
	for _, listener := range executor.testListeners {
		listener.AfterTest(r)
	}
}

func (executor executorImpl) executeSingleTest(t test.Test, device devices.Device, testPackage string, testRunner string) (string, error, time.Duration) {
	start := time.Now()
	testOutput, err := executor.adb.ExecuteTest(testPackage, testRunner, t.FullName(), device.Serial)
	duration := time.Since(start)
	return testOutput, err, duration
}
