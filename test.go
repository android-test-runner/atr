package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/ybonjour/atr/aapt"
	"github.com/ybonjour/atr/apk"
	"github.com/ybonjour/atr/device"
	"github.com/ybonjour/atr/test"
)

var testCommand = cli.Command{
	Name:   "test",
	Usage:  "Execute an android instrumentation test",
	Action: testAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "apk, a",
			Usage: "APK under test",
		},
		cli.StringFlag{
			Name:  "testapk, ta",
			Usage: "APK containing instrumentation tests",
		},
		cli.StringSliceFlag{
			Name:  "test, t",
			Usage: "Test to run formatted as TestClass#test",
		},
		cli.StringFlag{
			Name:  "testfile, tf",
			Usage: "Path to a textfile defining the tests to execute separated by newlines",
		},
		cli.StringSliceFlag{
			Name:  "device, d",
			Usage: "Id of device on which the test shall run",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "build/atr",
			Usage: "Folder to write test output",
		},
	},
}

func testAction(c *cli.Context) error {
	apkPath := c.String("apk")
	apkUnderTest, apkGetError := apk.GetApk(apkPath)
	if apkGetError != nil {
		return cli.NewExitError(fmt.Sprintf("Could not get apk %v", apkPath), 1)
	}

	testApkPath := c.String("testapk")
	testApk, apkGetError := apk.GetApk(testApkPath)
	if apkGetError != nil {
		return cli.NewExitError(fmt.Sprintf("Could not get apk %v", testApkPath), 1)
	}

	allTests, testsError := allTests(c)
	if testsError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid tests: %v", testsError), 1)
	}

	devices, devicesError := devices(c)
	if devicesError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid devices: %v", devicesError), 1)
	}

	testRunner, testRunnerError := aapt.TestRunner(testApk.Path)
	if testRunnerError != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid test runner: %v", testRunnerError), 1)
	}

	config := test.TestConfig{
		Apk:          apkUnderTest,
		TestApk:      testApk,
		TestRunner:   testRunner,
		Tests:        allTests,
		OutputFolder: c.String("output"),
	}

	return test.ExecuteTests(config, devices)
}

func devices(c *cli.Context) ([]device.Device, error) {
	deviceFlag := c.StringSlice("device")
	// ensure no filter (=> all connected devices) if no devices provided
	var d []string
	if len(deviceFlag) > 0 {
		d = deviceFlag
	}
	return device.ConnectedDevices(d)
}

func allTests(c *cli.Context) ([]test.Test, error) {
	var tests []test.Test
	testFile := c.String("testfile")
	if testFile != "" {
		testsFromFile, err := test.ParseTestsFromFile(testFile)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testsFromFile...)
	}

	tests = append(tests, test.ParseTests(c.StringSlice("test"))...)

	return tests, nil
}
