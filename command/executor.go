package command

import (
	"bytes"
	"os/exec"
)

type ExecutionResult struct {
	StdOut string
	StdErr string
	Error  error
}

type Executor interface {
	Execute(cmd *exec.Cmd) ExecutionResult
	ExecuteInBackground(cmd *exec.Cmd) (int, error)
}

type executorImpl struct{}

func NewExecutor() Executor {
	return executorImpl{}
}

func (executorImpl) Execute(cmd *exec.Cmd) ExecutionResult {
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
