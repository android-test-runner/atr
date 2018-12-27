package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"github.com/ybonjour/atr/apk"
	"github.com/ybonjour/atr/device"
	"os"
)

var testCommand = cli.Command{
	Name:   "test",
	Usage:  "Execute an android instrumentation test",
	Action: test,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "apk, a",
			Usage: "APK under test",
		},
		cli.StringFlag{
			Name:  "testapk, ta",
			Usage: "APK containing instrumentation tests",
		},
		cli.StringFlag{
			Name:  "testrunner, tr",
			Value: "android.support.test.runner.AndroidJUnitRunner",
			Usage: "Test Runner to run instrumentation tests",
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
	},
}

func test(c *cli.Context) error {
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

	config := TestConfig{
		Apk:        apkUnderTest,
		TestApk:    testApk,
		TestRunner: c.String("testrunner"),
		Tests:      allTests,
	}

	return ExecuteTests(config, devices)
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

func allTests(c *cli.Context) ([]string, error) {
	var tests []string

	if path := c.String("testfile"); path != "" {
		testsFromFile, err := testsFromFile(path)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testsFromFile...)
	}

	testsFromFlags := c.StringSlice("test")
	tests = append(tests, testsFromFlags...)

	return tests, nil
}

func testsFromFile(path string) ([]string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var tests []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tests = append(tests, scanner.Text())
	}

	return tests, scanner.Err()
}
