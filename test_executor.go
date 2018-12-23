package main

import "fmt"

type TestConfig struct {
	Apk     string
	TestApk string
}

func executeTests(config TestConfig) error {
	fmt.Printf("%v\n", config.Apk)
	fmt.Printf("%v\n", config.TestApk)

	return nil
}
