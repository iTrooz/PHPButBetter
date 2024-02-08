package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type CmdError struct {
	Cmd    []string
	Stderr []byte
	Err    error
}

func (err *CmdError) Error() string {
	return fmt.Sprintf("Command '%s' failed: %v\nStderr: %v", strings.Join(err.Cmd, " "), err.Err, string(err.Stderr))
}

func RunCmd(name string, args ...string) (*bytes.Buffer, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, &CmdError{
			Cmd:    append([]string{name}, args...),
			Stderr: stderr.Bytes(),
			Err:    err,
		}
	}
	return &stdout, nil
}
