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

type Parser interface {
	ParseFromFile(path string) ([]Test, error)
	Parse(testsUnparsed []string) []Test
}

type parserImpl struct{}

func NewParser() Parser {
	return parserImpl{}
}

func (parser parserImpl) ParseFromFile(path string) ([]Test, error) {
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

	return parser.Parse(tests), scanner.Err()
}

func (parserImpl) Parse(testsUnparsed []string) []Test {
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
