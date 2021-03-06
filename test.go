package main

import (
	"fmt"
	"github.com/android-test-runner/atr/aapt"
	"github.com/android-test-runner/atr/apks"
	"github.com/android-test-runner/atr/console"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/junit_xml"
	"github.com/android-test-runner/atr/logcat"
	"github.com/android-test-runner/atr/logging"
	"github.com/android-test-runner/atr/output"
	"github.com/android-test-runner/atr/screen_recorder"
	"github.com/android-test-runner/atr/test"
	"github.com/android-test-runner/atr/test_executor"
	"github.com/android-test-runner/atr/test_listener"
	"github.com/urfave/cli/altsrc"
	"gopkg.in/urfave/cli.v1"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "load, l",
		Usage: "Specify file to load flag values from.",
	},
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "apk",
		Usage: "APK under test.",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "testapk",
		Usage: "APK containing instrumentation tests.",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:  "test",
		Usage: "Test to run formatted as TestClass#testMethod.",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "testfile",
		Usage: "Path to a text file defining the tests to execute separated by newlines.",
	}),
	altsrc.NewStringSliceFlag(cli.StringSliceFlag{
		Name:  "device",
		Usage: "Serial number of device on which the test shall run.",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:  "output",
		Value: "build/atr",
		Usage: "Folder to write test output to.",
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
		Usage: "Record test results in JUnit XML report.",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "disableanimations",
		Usage: "Disable animations before test execution. This is recommended for Espresso tests.",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "verbose",
		Usage: "Enables debug messages",
	}),
}

var testCommand = cli.Command{
	Name:   "test",
	Usage:  "Execute an android instrumentation test",
	Before: readYaml,
	Action: testAction,
	Flags:  flags,
}

func readYaml(c *cli.Context) error {
	if c.String("load") == "" {
		return nil
	}

	err := altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load"))(c)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not read YAML configuration file '%v': \n%v\n", c.String("load"), err)
		return cli.NewExitError(errorMessage, 1)
	}
	return nil
}

func testAction(c *cli.Context) error {
	if c.Bool("verbose") {
		logging.SetLogLevel(logging.Debug)
	}

	apkPath := c.String("apk")
	apkUnderTest, apkGetError := apks.New().GetApk(apkPath)
	if apkGetError != nil {
		return cli.NewExitError(fmt.Sprintf("Could not get apk '%v'", apkPath), 1)
	}

	testApkPath := c.String("testapk")
	testApk, apkGetError := apks.New().GetApk(testApkPath)
	if apkGetError != nil {
		return cli.NewExitError(fmt.Sprintf("Could not get testapk '%v'", testApkPath), 1)
	}

	allTests, testsError := allTests(c)
	if testsError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid tests: '%v'", testsError), 1)
	}

	configDevices, devicesError := getDevices(c)
	if devicesError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid devices: '%v'", devicesError), 1)
	}

	fmt.Printf("Executing %v Tests on the following devices: '%v'\n", len(allTests), configDevices)

	testRunner, testRunnerError := aapt.New().TestRunner(testApk.Path)
	if testRunnerError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid test runner: %v", testRunnerError), 1)
	}

	config := test_executor.Config{
		Apk:               apkUnderTest,
		TestApk:           testApk,
		TestRunner:        testRunner,
		Tests:             allTests,
		OutputFolder:      c.String("output"),
		DisableAnimations: c.Bool("disableanimations"),
	}
	writer := output.NewWriter(c.String("output"))

	testListenersFactory := testListenerFactory{context: c, apk: apkUnderTest, writer: writer}
	testExecutionError := test_executor.NewExecutor(writer, testListenersFactory).Execute(config, configDevices)
	if testExecutionError != nil {
		return cli.NewExitError(fmt.Sprintf("Test execution errored:\n %v", testExecutionError), 1)
	}

	return nil
}

type testListenerFactory struct {
	context *cli.Context
	apk     apks.Apk
	writer  output.Writer
}

func (f testListenerFactory) ForDevice(device devices.Device) []test_listener.TestListener {
	var listeners []test_listener.TestListener

	listeners = append(listeners, console.NewTestListener(device))

	if f.context.Bool("recordjunit") {
		listeners = append(listeners, junit_xml.NewTestListener(device, f.writer, f.apk))
	}

	if f.context.Bool("recordlogcat") {
		listeners = append(listeners, logcat.NewTestListener(device, f.writer))
	}

	if f.context.Bool("recordscreen") {
		listeners = append(listeners, screen_recorder.NewTestListener(device, f.writer))
	}
	return listeners
}

func getDevices(c *cli.Context) ([]devices.Device, error) {
	deviceDefinitions := c.StringSlice("device")
	d := devices.New()
	if len(deviceDefinitions) > 0 {
		return d.ParseConnectedDevices(deviceDefinitions)
	} else {
		return d.AllConnectedDevices()
	}
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
