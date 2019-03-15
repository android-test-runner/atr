package test_listener

import (
	"github.com/android-test-runner/atr/result"
	"github.com/android-test-runner/atr/test"
)

type TestListener interface {
	BeforeTestSuite()
	AfterTestSuite()
	BeforeTest(test test.Test)
	AfterTest(r result.Result) []result.Extra
}
