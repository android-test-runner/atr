package command

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type Executor interface {
	Execute(cmd *exec.Cmd) error
	ExecuteOutput(cmd *exec.Cmd) (string, error)
}

type executorImpl struct{}

func NewExecutor() Executor {
	return executorImpl{}
}

func (executor executorImpl) Execute(cmd *exec.Cmd) error {
	out, executeError := executor.ExecuteOutput(cmd)

	if executeError != nil {
		return executeError
	}

	fmt.Printf("%v", out)
	return nil
}

func (executor executorImpl) ExecuteOutput(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	cmd.Stdout = &out

	var err bytes.Buffer
	cmd.Stderr = &err

	runError := cmd.Run()
	if runError != nil {
		log.Printf("%v", err.String())
		return "", runError
	}

	return out.String(), nil
}
