package test_executor

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/logcat"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/screen_recorder"
	"github.com/ybonjour/atr/test"
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
	installer             Installer
	resultParser          result.Parser
	adb                   adb.Adb
	logcatFactory         logcat.Factory
	screenRecorderFactory screen_recorder.Factory
}

func NewExecutor(writer output.Writer) Executor {
	return executorImpl{
		installer:             NewInstaller(),
		resultParser:          result.NewParser(),
		adb:                   adb.New(),
		logcatFactory:         logcat.NewFactory(writer),
		screenRecorderFactory: screen_recorder.NewFactory(writer),
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
	err := executor.reinstallApks(config, device)
	if err != nil {
		return nil, err
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
	deviceLogcat := executor.logcatFactory.ForDevice(device)
	screenRecorder := executor.screenRecorderFactory.ForDevice(device)
	for _, t := range testConfig.Tests {
		executor.beforeTest(t, deviceLogcat, screenRecorder)
		testOutput, errTest, duration := executor.executeSingleTest(t, device, testConfig.TestApk.PackageName, testConfig.TestRunner)
		executor.afterTest(t, deviceLogcat, screenRecorder)

		results = append(results, executor.resultParser.ParseFromOutput(t, errTest, testOutput, duration))
	}
	return results
}

func (executor executorImpl) beforeTest(t test.Test, logcat logcat.Logcat, recorder screen_recorder.ScreenRecorder) {
	errStartLogcat := logcat.StartRecording(t)
	if errStartLogcat != nil {
		fmt.Printf("Could not clear logcat: '%v'\n", errStartLogcat)
	}
	errStartScreenRecording := recorder.StartRecording(t)
	if errStartScreenRecording != nil {
		fmt.Printf("Could not start screen recording: '%v'\n", errStartScreenRecording)
	}
}

func (executor executorImpl) afterTest(t test.Test, logcat logcat.Logcat, recorder screen_recorder.ScreenRecorder) {
	errStopLogcat := logcat.StopRecording(t)
	if errStopLogcat != nil {
		fmt.Printf("Could not save logcat: '%v'\n", errStopLogcat)
	}
	errStopScreenRecording := recorder.StopRecording(t)
	if errStopScreenRecording != nil {
		fmt.Printf("Could not save screen recording: '%v'\n", errStopScreenRecording)
	}
}

func (executor executorImpl) executeSingleTest(t test.Test, device devices.Device, testPackage string, testRunner string) (string, error, time.Duration) {
	start := time.Now()
	testOutput, err := executor.adb.ExecuteTest(testPackage, testRunner, t.FullName(), device.Serial)
	duration := time.Since(start)
	return testOutput, err, duration
}
