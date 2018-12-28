package test_executor

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/mock_adb"
	"github.com/ybonjour/atr/mock_result"
	"github.com/ybonjour/atr/mock_test_executor"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"testing"
)

func TestExecute(t *testing.T) {
	targetTest := test.Test{Class: "TestClass", Method: "testMethod"}
	config := Config{
		Apk:        apks.Apk{},
		TestApk:    apks.Apk{PackageName: "testPackageName"},
		Tests:      []test.Test{targetTest},
		TestRunner: "testRunner",
	}
	testOutput := "testOutput"
	testResult := result.Result{}
	device := devices.Device{Serial: "abcd"}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstaller := mock_test_executor.NewMockInstaller(ctrl)
	mockInstaller.EXPECT().Reinstall(config.Apk, device).Return(nil)
	mockInstaller.EXPECT().Reinstall(config.TestApk, device).Return(nil)
	mockAdb := mock_adb.NewMockAdb(ctrl)
	mockAdb.
		EXPECT().
		ExecuteTest(config.TestApk.PackageName, config.TestRunner, targetTest.FullName(), device.Serial).
		Return(testOutput, nil)
	mockResultParser := mock_result.NewMockResultParser(ctrl)
	mockResultParser.EXPECT().ParseFromOutput(targetTest, nil, testOutput).Return(testResult)
	executor := executorImpl{
		installer:    mockInstaller,
		adb:          mockAdb,
		resultParser: mockResultParser,
	}

	results, err := executor.Execute(config, []devices.Device{device})

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
	expectedResults := []result.Result{testResult}
	if !AreEqualResults(results[device], expectedResults) {
		t.Error(fmt.Sprintf("Expected results '%v' but got '%v'", expectedResults, results[device]))
	}
}

func AreEqualResults(slice1, slice2 []result.Result) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}
