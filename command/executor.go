package command

import (
	"bytes"
	"log"
	"os/exec"
)

type ExecutionResult struct {
	StdOut string
	StdErr string
	Error  error
}

type Executor interface {
	ExecuteResult(cmd *exec.Cmd) ExecutionResult
	ExecuteOutput(cmd *exec.Cmd) (string, error)
	ExecuteInBackground(cmd *exec.Cmd) (int, error)
}

type executorImpl struct{}

func NewExecutor() Executor {
	return executorImpl{}
}

func (executorImpl) ExecuteResult(cmd *exec.Cmd) ExecutionResult {
	var out bytes.Buffer
	cmd.Stdout = &out

	var err bytes.Buffer
	cmd.Stderr = &err

	runError := cmd.Run()

	return ExecutionResult{
		StdOut: out.String(),
		StdErr: err.String(),
		Error:  runError,
	}
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
