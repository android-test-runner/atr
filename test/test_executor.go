package test

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apk"
	"github.com/ybonjour/atr/device"
	"regexp"
	"strings"
)

type TestConfig struct {
	Apk          *apk.Apk
	TestApk      *apk.Apk
	TestRunner   string
	Tests        []Test
	OutputFolder string
}

type TestResult struct {
	Test      Test
	HasPassed bool
	Output    string
}

func ExecuteTests(testConfig TestConfig, devices []device.Device) error {
	for _, d := range devices {
		apkInstallError := reinstall(testConfig.Apk, d)
		if apkInstallError != nil {
			return apkInstallError
		}
		testApkInstallError := reinstall(testConfig.TestApk, d)
		if testApkInstallError != nil {
			return testApkInstallError
		}
		testResults := executeTests(testConfig, d)

		fmt.Printf("Results %v\n", testResults)
	}
	return nil
}

func executeTests(testConfig TestConfig, device device.Device) []TestResult {
	var results []TestResult
	for _, t := range testConfig.Tests {
		output, testError := adb.ExecuteTest(testConfig.TestApk.PackageName, testConfig.TestRunner, FullName(t), device.Serial)

		hasPassed := testError == nil && testSuccessful(output)
		result := TestResult{
			Test:      t,
			HasPassed: hasPassed,
			Output:    output,
		}
		results = append(results, result)
	}

	return results
}

func testSuccessful(output string) bool {
	// A test was successful if we find "OK (1 test)" in the output
	// This is needed because the am process does not fail if the test fails.
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		regexOk := regexp.MustCompile(`^OK \(1 test\)$`)
		if regexOk.MatchString(line) {
			return true
		}
	}

	return false
}

func reinstall(apk *apk.Apk, device device.Device) error {
	apkUninstallError := adb.Uninstall(apk.PackageName, device.Serial)
	if apkUninstallError != nil {
		fmt.Println("Could not uninstall apk. Try to install it anyways.")
	}

	apkInstallError := adb.Install(apk.Path, device.Serial)
	if apkInstallError != nil {
		return apkInstallError
	}

	return nil
}
