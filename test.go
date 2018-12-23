package main

import (
	"fmt"
	"github.com/urfave/cli"
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
	},
}

func test(c *cli.Context) error {
	fmt.Printf("%v\n", c.String("apk"))
	fmt.Printf("%v\n", c.String("testapk"))
	return nil
}
