package command

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type executionResult struct {
	StdOut string
	StdErr string
	Error  error
}

type Executor interface {
	Execute(cmd *exec.Cmd) error
	ExecuteResult(cmd *exec.Cmd) executionResult
	ExecuteOutput(cmd *exec.Cmd) (string, error)
	ExecuteInBackground(cmd *exec.Cmd) (int, error)
}

type executorImpl struct{}

func NewExecutor() Executor {
	return executorImpl{}
}

func (executorImpl) ExecuteResult(cmd *exec.Cmd) executionResult {
	var out bytes.Buffer
	cmd.Stdout = &out

	var err bytes.Buffer
	cmd.Stderr = &err

	runError := cmd.Run()

	return executionResult{
		StdOut: out.String(),
		StdErr: err.String(),
		Error:  runError,
	}
}

func (executor executorImpl) Execute(cmd *exec.Cmd) error {
	out, executeError := executor.ExecuteOutput(cmd)

	if executeError != nil {
		return executeError
	}

	fmt.Printf("%v", out)
	return nil
}

func (executor executorImpl) ExecuteInBackground(cmd *exec.Cmd) (int, error) {
	err := cmd.Start()
	return cmd.Process.Pid, err
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
