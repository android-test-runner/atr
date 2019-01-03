# Android Test Runner (ATR)
A utility to run Android instrumentation tests from command line.

## Vision
Android Test Runner shall support he following functionality

* Run all tests in apk-test against apk-under-test
`atr test --apk=under-test.apk --testapk=test.apk` 

* Run a single test
`atr test --test="MyTest#test"`

* Run multiple tests
`atr test --test="MyTest#test1" --test="MyTest#test2"`

* Record screen during test execution
`atr test --record-screen=true`

* Record logcat during test execution
`atr test --record-logcat=true`

* List connected devices
`atr devices`

* Run tests on a specific device
atr test `--device="abcjekrjkdfj43r"`

* `.atr` file to set default options like apk, testapk, filter device, record-screen etc. for project or user

* adb retries to stabalize adb


## Dependencies
* adb
* aapt


## Notes / Links
* [Run tests through adb](https://developer.android.com/studio/test/command-line) `adb shell am instrument -w <test_package_name>/<runner_class>`
* [Find Packgaename](https://stackoverflow.com/questions/6289149/read-the-package-name-of-an-android-apk):
`aapt dump badging <path-to-apk> | grep package:\ name`
* [Parse tests in APK (JVM required)](https://github.com/linkedin/dex-test-parser): Download jar from https://dl.bintray.com/linkedin/maven/com/linkedin/dextestparser/parser/2.0.1/parser-2.0.1-all.jar
`java -jar parser-2.0.1-all.jar test.apk output-file`


