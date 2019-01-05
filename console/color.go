package console

import "fmt"

const (
	Red     = "\033[0;31m"
	Green   = "\033[0;32m"
	noColor = "\033[0m"
)

func Color(input string, color string) string {
	return fmt.Sprintf("%v%v%v", color, input, noColor)
}
