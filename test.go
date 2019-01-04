package main

import (
	"fmt"
	"github.com/urfave/cli/altsrc"
	"github.com/ybonjour/atr/aapt"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/files"
	"github.com/ybonjour/atr/junit_xml"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_executor"
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

	resultsByDevice, testExecutionError := test_executor.NewExecutor(writer).Execute(config, configDevices)
	if testExecutionError != nil {
		return cli.NewExitError(fmt.Sprintf("Test execution errored: '%v'", testExecutionError), 1)
	}
	junitXmlFiles := toJunitXmlFiles(resultsByDevice, apkUnderTest)

	err := writer.Write(junitXmlFiles)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error while writing junit results '%v'", err), 1)
	}

	return nil
}

func toJunitXmlFiles(resultsByDevice map[devices.Device][]result.Result, apk apks.Apk) map[devices.Device][]files.File {
	xmlFiles := map[devices.Device][]files.File{}
	for device, results := range resultsByDevice {
		formatter := junit_xml.NewFormatter()
		xmlFile, err := formatter.Format(results, apk)
		if err != nil {
			fmt.Printf("Error while formatting test results for device '%v': '%v'. Ignoring results for this device and continue with next device.\n", err)
			continue
		}
		xmlFiles[device] = []files.File{xmlFile}
	}
	return xmlFiles
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
