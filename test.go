package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/ybonjour/atr/apk"
	"github.com/ybonjour/atr/device"
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

	config := TestConfig{
		Apk:        apkUnderTest,
		TestApk:    testApk,
		TestRunner: c.String("testrunner"),
		Tests:      c.StringSlice("test"),
	}

	devices := []device.Device{}
	deviceSerials := c.StringSlice("device")
	if len(deviceSerials) == 0 {
		allDevices, err := device.AllConnectedDevices()
		if err != nil {
			return err
		}

		devices = allDevices

	} else {
		devices = device.FromSerials(deviceSerials)
	}

	return ExecuteTests(config, devices)
}
