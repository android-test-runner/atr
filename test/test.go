package test

import "fmt"

type Test struct {
	Class  string
	Method string
}

func (t Test) FullName() string {
	return fmt.Sprintf("%v#%v", t.Class, t.Method)
}
