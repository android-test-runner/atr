package main

import (
	"fmt"
	"github.com/urfave/cli/altsrc"
	"github.com/ybonjour/atr/aapt"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/junit_xml"
	"github.com/ybonjour/atr/logcat"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/screen_recorder"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_executor"
	"github.com/ybonjour/atr/test_listener"
	"gopkg.in/urfave/cli.v1"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "load, l",
		Usage: "Specify to file to load flags from",
	},
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "apk",
		Usage: "APK under test",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "testapk",
		Usage: "APK containing instrumentation tests",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:  "test",
		Usage: "Test to run formatted as TestClass#test",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "testfile",
		Usage: "Path to a textfile defining the tests to execute separated by newlines",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:  "device",
		Usage: "Id of device on which the test shall run",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "output",
		Value: "build/atr",
		Usage: "Folder to write test output",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "recordscreen",
		Usage: "Record screen for failed tests.",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "recordlogcat",
		Usage: "Record logcat for failed tests.",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "recordjunit",
		Usage: "Records test results in JUnit XML report",
	}),
}

var testCommand = cli.Command{
	Name:   "test",
	Usage:  "Execute an android instrumentation test",
	Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load")),
	Action: testAction,
	Flags:  flags,
}

func testAction(c *cli.Context) error {
	apkPath := c.String("apk")
	apkUnderTest, apkGetError := apks.New().GetApk(apkPath)
	if apkGetError != nil {
		return cli.NewExitError(fmt.Sprintf("Could not get apk %v", apkPath), 1)
	}

	testApkPath := c.String("testapk")
	testApk, apkGetError := apks.New().GetApk(testApkPath)
	if apkGetError != nil {
		return cli.NewExitError(fmt.Sprintf("Could not get apk %v", testApkPath), 1)
	}

	allTests, testsError := allTests(c)
	if testsError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid tests: %v", testsError), 1)
	}

	configDevices, devicesError := getDevices(c)
	if devicesError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid devices: %v", devicesError), 1)
	}

	fmt.Printf("Executing %v Tests on the following devices: '%v'", len(allTests), configDevices)

	testRunner, testRunnerError := aapt.New().TestRunner(testApk.Path)
	if testRunnerError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid test runner: %v", testRunnerError), 1)
	}

	config := test_executor.Config{
		Apk:          apkUnderTest,
		TestApk:      testApk,
		TestRunner:   testRunner,
		Tests:        allTests,
		OutputFolder: c.String("output"),
	}
	writer := output.NewWriter(c.String("output"))

	testListeners := getTestListeners(c, apkUnderTest, writer)
	testExecutionError := test_executor.NewExecutor(writer, testListeners).Execute(config, configDevices)
	if testExecutionError != nil {
		return cli.NewExitError(fmt.Sprintf("Test execution errored: '%v'", testExecutionError), 1)
	}

	return nil
}

func getTestListeners(c *cli.Context, apk apks.Apk, writer output.Writer) []test_listener.TestListener {
	var listeners []test_listener.TestListener

	if c.Bool("recordjunit") {
		listeners = append(listeners, junit_xml.NewTestListener(writer, apk))
	}

	if c.Bool("recordlogcat") {
		listeners = append(listeners, logcat.NewTestListener(writer))
	}

	if c.Bool("recordscreen") {
		listeners = append(listeners, screen_recorder.NewTestListener(writer))
	}

	return listeners
}

func getDevices(c *cli.Context) ([]devices.Device, error) {
	deviceFlag := c.StringSlice("device")
	// ensure no filter (=> all connected devices) if no devices provided
	var d []string
	if len(deviceFlag) > 0 {
		d = deviceFlag
	}
	return devices.New().ConnectedDevices(d)
}

func allTests(c *cli.Context) ([]test.Test, error) {
	var tests []test.Test
	testFile := c.String("testfile")
	if testFile != "" {
		testsFromFile, err := test.NewParser().ParseFromFile(testFile)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testsFromFile...)
	}

	tests = append(tests, test.NewParser().Parse(c.StringSlice("test"))...)

	return tests, nil
}
