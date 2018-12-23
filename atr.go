package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var version string //needs to be set when compiling using ldflags (-ldflags "-X main.version=42")

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
		cli.StringSliceFlag{
			Name:  "test, t",
			Usage: "Test to execute",
		},
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "atr"
	app.Usage = "Android Test Runner"
	app.Version = version
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		testCommand,
	}

	app.Run(os.Args)
}

func test(c *cli.Context) error {
	fmt.Printf("%v\n", c.String("apk"))
	fmt.Printf("%v\n", c.String("testapk"))
	fmt.Printf("%v\n", c.String("test"))
	return nil
}
