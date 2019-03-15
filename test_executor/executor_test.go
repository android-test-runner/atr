package test_executor

import (
	"fmt"
	"github.com/android-test-runner/atr/apks"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
	"github.com/android-test-runner/atr/mock_adb"
	"github.com/android-test-runner/atr/mock_output"
	"github.com/android-test-runner/atr/mock_result"
	"github.com/android-test-runner/atr/mock_test_executor"
	"github.com/android-test-runner/atr/mock_test_listener"
	"github.com/android-test-runner/atr/result"
	"github.com/android-test-runner/atr/test"
	"github.com/android-test-runner/atr/test_listener"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestExecute(t *testing.T) {
	targetTest := test.Test{Class: "TestClass", Method: "testMethod"}
	config := Config{
		TestApk:           apks.Apk{PackageName: "testPackageName"},
		Tests:             []test.Test{targetTest},
		TestRunner:        "testRunner",
		DisableAnimations: true,
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
	mockAdb.EXPECT().DisableAnimations(device.Serial).Return(nil)
	mockResultParser := mock_result.NewMockParser(ctrl)
	mockResultParser.EXPECT().ParseFromOutput(gomock.Eq(targetTest), gomock.Eq(nil), gomock.Eq(testOutput), gomock.Any()).Return(testResult)
	mockJsonFormatter := mock_result.NewMockJsonFormatter(ctrl)
	mockHtmlFormatter := mock_result.NewMockHtmlFormatter(ctrl)
	mockWriter := mock_output.NewMockWriter(ctrl)
	givenDeviceDirectoryCanBeRemoved(device, mockWriter)
	givenJsonFileCanBeWritten(mockJsonFormatter, mockWriter)
	givenHtmlFileCanBeWritten(mockHtmlFormatter, mockWriter)
	mockTestListenerFactory := mock_test_listener.NewMockFactory(ctrl)
	givenNoTestListeners(device, mockTestListenerFactory)
	executor := executorImpl{
		installer:            mockInstaller,
		adb:                  mockAdb,
		resultParser:         mockResultParser,
		testListenersFactory: mockTestListenerFactory,
		jsonFormatter:        mockJsonFormatter,
		htmlFormatter:        mockHtmlFormatter,
		writer:               mockWriter,
	}

	err := executor.Execute(config, []devices.Device{device})

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
}

func TestExecuteMultipleTests(t *testing.T) {
	test1 := test.Test{Class: "TestClass", Method: "testMethod"}
	test2 := test.Test{Class: "TestClass", Method: "testMethod1"}
	testResult1 := result.Result{}
	testResult2 := result.Result{}
	device := devices.Device{Serial: "abcd"}
	config := Config{
		Tests: []test.Test{test1, test2},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstaller := mock_test_executor.NewMockInstaller(ctrl)
	mockAdb := mock_adb.NewMockAdb(ctrl)
	mockResultParser := mock_result.NewMockParser(ctrl)
	givenAllApksInstalledSuccessfully(mockInstaller, 1)
	givenTestOnDeviceReturns(test1, device, testResult1, mockAdb, mockResultParser)
	givenTestOnDeviceReturns(test2, device, testResult2, mockAdb, mockResultParser)
	mockJsonFormatter := mock_result.NewMockJsonFormatter(ctrl)
	mockHtmlFormatter := mock_result.NewMockHtmlFormatter(ctrl)
	mockWriter := mock_output.NewMockWriter(ctrl)
	givenDeviceDirectoryCanBeRemoved(device, mockWriter)
	givenJsonFileCanBeWritten(mockJsonFormatter, mockWriter)
	givenHtmlFileCanBeWritten(mockHtmlFormatter, mockWriter)
	mockTestListenerFactory := mock_test_listener.NewMockFactory(ctrl)
	givenNoTestListeners(device, mockTestListenerFactory)
	executor := executorImpl{
		installer:            mockInstaller,
		adb:                  mockAdb,
		resultParser:         mockResultParser,
		testListenersFactory: mockTestListenerFactory,
		jsonFormatter:        mockJsonFormatter,
		htmlFormatter:        mockHtmlFormatter,
		writer:               mockWriter,
	}

	err := executor.Execute(config, []devices.Device{device})

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
}

func TestExecuteMultipleDevices(t *testing.T) {
	targetTest := test.Test{Class: "TestClass", Method: "testMethod"}
	testResult1 := result.Result{}
	testResult2 := result.Result{}
	device1 := devices.Device{Serial: "abcd"}
	device2 := devices.Device{Serial: "efgh"}
	config := Config{
		Tests: []test.Test{targetTest},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstaller := mock_test_executor.NewMockInstaller(ctrl)
	mockAdb := mock_adb.NewMockAdb(ctrl)
	mockResultParser := mock_result.NewMockParser(ctrl)
	givenAllApksInstalledSuccessfully(mockInstaller, 2)
	givenTestOnDeviceReturns(targetTest, device1, testResult1, mockAdb, mockResultParser)
	givenTestOnDeviceReturns(targetTest, device2, testResult2, mockAdb, mockResultParser)
	mockJsonFormatter := mock_result.NewMockJsonFormatter(ctrl)
	mockHtmlFormatter := mock_result.NewMockHtmlFormatter(ctrl)
	mockWriter := mock_output.NewMockWriter(ctrl)
	givenDeviceDirectoryCanBeRemoved(device1, mockWriter)
	givenDeviceDirectoryCanBeRemoved(device2, mockWriter)
	givenJsonFileCanBeWritten(mockJsonFormatter, mockWriter)
	givenHtmlFileCanBeWritten(mockHtmlFormatter, mockWriter)
	mockTestListenerFactory := mock_test_listener.NewMockFactory(ctrl)
	givenNoTestListeners(device1, mockTestListenerFactory)
	givenNoTestListeners(device2, mockTestListenerFactory)
	executor := executorImpl{
		installer:            mockInstaller,
		adb:                  mockAdb,
		resultParser:         mockResultParser,
		testListenersFactory: mockTestListenerFactory,
		jsonFormatter:        mockJsonFormatter,
		htmlFormatter:        mockHtmlFormatter,
		writer:               mockWriter,
	}

	err := executor.Execute(config, []devices.Device{device1, device2})

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
}

func TestExecuteCallsTestListener(t *testing.T) {
	targetTest := test.Test{Class: "TestClass", Method: "testMethod"}
	config := Config{
		TestApk:    apks.Apk{PackageName: "testPackageName"},
		Tests:      []test.Test{targetTest},
		TestRunner: "testRunner",
	}
	testResult := result.Result{}
	device := devices.Device{Serial: "abcd"}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstaller := mock_test_executor.NewMockInstaller(ctrl)
	givenAllApksInstalledSuccessfully(mockInstaller, 1)
	mockAdb := mock_adb.NewMockAdb(ctrl)
	mockResultParser := mock_result.NewMockParser(ctrl)
	givenTestOnDeviceReturns(targetTest, device, testResult, mockAdb, mockResultParser)
	mockJsonFormatter := mock_result.NewMockJsonFormatter(ctrl)
	mockHtmlFormatter := mock_result.NewMockHtmlFormatter(ctrl)
	mockWriter := mock_output.NewMockWriter(ctrl)
	givenDeviceDirectoryCanBeRemoved(device, mockWriter)
	givenJsonFileCanBeWritten(mockJsonFormatter, mockWriter)
	givenHtmlFileCanBeWritten(mockHtmlFormatter, mockWriter)
	testListener := mock_test_listener.NewMockTestListener(ctrl)
	testListener.EXPECT().BeforeTestSuite()
	testListener.EXPECT().BeforeTest(targetTest)
	testListener.EXPECT().AfterTest(testResult)
	testListener.EXPECT().AfterTestSuite()
	testListenerFactory := mock_test_listener.NewMockFactory(ctrl)
	testListenerFactory.EXPECT().ForDevice(device).Return([]test_listener.TestListener{testListener})
	executor := executorImpl{
		installer:            mockInstaller,
		adb:                  mockAdb,
		resultParser:         mockResultParser,
		testListenersFactory: testListenerFactory,
		jsonFormatter:        mockJsonFormatter,
		htmlFormatter:        mockHtmlFormatter,
		writer:               mockWriter,
	}

	err := executor.Execute(config, []devices.Device{device})

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
}

func givenNoTestListeners(device devices.Device, testListenerFactory *mock_test_listener.MockFactory) {
	testListenerFactory.EXPECT().ForDevice(device).Return([]test_listener.TestListener{})
}

func givenDeviceDirectoryCanBeRemoved(device devices.Device, mockWriter *mock_output.MockWriter) {
	mockWriter.EXPECT().RemoveDeviceDirectory(device).Return(nil)
}

func givenAllApksInstalledSuccessfully(mockInstaller *mock_test_executor.MockInstaller, numDevices int) {
	mockInstaller.EXPECT().Reinstall(gomock.Any(), gomock.Any()).Return(nil).Times(numDevices * 2)
}

func givenTestOnDeviceReturns(t test.Test, d devices.Device, r result.Result, mockAdb *mock_adb.MockAdb, mockResultParser *mock_result.MockParser) {
	testOutput := t.FullName()
	mockAdb.
		EXPECT().
		ExecuteTest(gomock.Any(), gomock.Any(), gomock.Eq(t.FullName()), gomock.Eq(d.Serial)).
		Return(testOutput, nil)

	mockResultParser.
		EXPECT().
		ParseFromOutput(gomock.Eq(t), gomock.Eq(nil), gomock.Eq(testOutput), gomock.Any()).
		Return(r)
}

func givenHtmlFileCanBeWritten(mockHtmlFormatter *mock_result.MockHtmlFormatter, mockWriter *mock_output.MockWriter) {
	f := files.File{}
	mockHtmlFormatter.EXPECT().FormatResults(gomock.Any()).Return(f, nil)
	mockWriter.EXPECT().WriteFileToRoot(f).Return("", nil)
}

func givenJsonFileCanBeWritten(mockJsonFormatter *mock_result.MockJsonFormatter, mockWriter *mock_output.MockWriter) {
	f := files.File{}
	mockJsonFormatter.EXPECT().FormatResults(gomock.Any()).Return(f, nil)
	mockWriter.EXPECT().WriteFileToRoot(f).Return("", nil)
}
