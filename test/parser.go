package test

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Test struct {
	Class  string
	Method string
}

func ParseTestsFromFile(path string) ([]Test, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var tests []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tests = append(tests, scanner.Text())
	}

	return ParseTests(tests), scanner.Err()
}

func ParseTests(testsUnparsed []string) []Test {
	var tests []Test
	for _, testUnparsed := range testsUnparsed {
		tests = append(tests, parseTest(testUnparsed))
	}

	return tests
}

func FullName(test Test) string {
	return fmt.Sprintf("%v#%v", test.Class, test.Method)
}

func parseTest(testUnparsed string) Test {
	tokens := strings.Split(testUnparsed, "#")
	return Test{
		Class:  tokens[0],
		Method: tokens[1],
	}
}
