package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/ybonjour/atr/apk"
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
			Value: "AndroidJUnitRunner",
			Usage: "Test Runner to run instrumentation tests",
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
	}
	return ExecuteTests(config)
}
