package test_executor

import (
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"sync"
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
	installer    Installer
	resultParser result.Parser
	adb          adb.Adb
}

func NewExecutor() Executor {
	return executorImpl{
		installer:    NewInstaller(),
		resultParser: result.NewParser(),
		adb:          adb.New(),
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
	for _, t := range testConfig.Tests {
		output, err := executor.adb.ExecuteTest(testConfig.TestApk.PackageName, testConfig.TestRunner, t.FullName(), device.Serial)
		results = append(results, executor.resultParser.ParseFromOutput(t, err, output))
	}

	return results
}
