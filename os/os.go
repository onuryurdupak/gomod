package os

import (
	"fmt"
	"os/exec"
)

// Executes input with bash returns stdout output.
func Exec(cmd string) (string, error) {
	executedCmd := exec.Command("/bin/bash", "-c", cmd)
	cmdOutput, err := executedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing input: '%s', err: %w", cmd, err)
	}

	return string(cmdOutput), nil
}
