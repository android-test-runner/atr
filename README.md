# Android Test Runner (ATR)
A utility to run Android instrumentation tests from command line.

## Examples
* Run tests listed in `AllTests.txt` in `test.apk` against `under-test.apk` on device with serial `abcdefg`:
`atr test --apk=under-test.apk --testapk=test.apk --testfile AllTests.txt --recordjunit --recordscreen --recordlogcat --device abcdefg` 

* Run tests `MyTestClass#myTest` in `test.apk` against `under-test.apk` on all connected devices
`atr test --apk=under-test.apk --testapk=test.apk --testfile AllTests.txt --recordjunit --recordscreen --recordlogcat`

* Run tests as specified in `atr.yaml` (You can provide flag values in this file using YAML)
`atr test --load atr.yaml`

## Dependencies
* adb
* aapt

## Notes / Links
* [Parse tests in APK (JVM required)](https://github.com/linkedin/dex-test-parser): Download jar from https://dl.bintray.com/linkedin/maven/com/linkedin/dextestparser/parser/2.0.1/parser-2.0.1-all.jar
`java -jar parser-2.0.1-all.jar test.apk output-file`


