package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

var version string //needs to be set when compiling using ldflags (-ldflags "-X main.version=42")

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
