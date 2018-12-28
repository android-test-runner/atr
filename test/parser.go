package test

import (
	"fmt"
	"github.com/ybonjour/atr/files"
	"strings"
)

type Test struct {
	Class  string
	Method string
}

func (t Test) FullName() string {
	return fmt.Sprintf("%v#%v", t.Class, t.Method)
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
	lines, err := files.New().ReadLines(path)
	if err != nil {
		return nil, err
	}
	return parser.Parse(lines), nil
}

func (parserImpl) Parse(testsUnparsed []string) []Test {
	var tests []Test
	for _, testUnparsed := range testsUnparsed {
		tests = append(tests, parseTest(testUnparsed))
	}

	return tests
}

func parseTest(testUnparsed string) Test {
	tokens := strings.Split(testUnparsed, "#")
	return Test{
		Class:  tokens[0],
		Method: tokens[1],
	}
}
