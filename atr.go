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
	Action: executeTest,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "apk-under-test",
			Usage: "The apk under test",
		},
		cli.StringFlag{
			Name:  "apk-test",
			Usage: "The test apk",
		},
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "atr"
	app.Usage = "Android Test Runner"
	app.Version = version
	app.Commands = []cli.Command{
		testCommand,
	}

	app.Run(os.Args)
}

func executeTest(c *cli.Context) error {
	apk := c.String("apk-under-test")
	fmt.Sprintf("Apk under test %v", apk)

	return nil
}
