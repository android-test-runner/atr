package command

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func Execute(cmd *exec.Cmd) error {
	out, executeError := ExecuteOutput(cmd)

	if executeError != nil {
		return executeError
	}

	fmt.Printf("%v", out)
	return nil
}

func ExecuteOutput(cmd *exec.Cmd) (string, error) {
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
