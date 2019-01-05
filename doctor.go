package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
)

var doctorCommand = cli.Command{
	Name:   "doctor",
	Usage:  "Helps with diagnosing atr problems",
	Action: doctorAction,
}

func doctorAction(c *cli.Context) {
	fmt.Println("Doctor")
}
