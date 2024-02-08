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

func RunCmd(cmd *exec.Cmd) (*bytes.Buffer, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, &CmdError{
			Cmd:    cmd.Args,
			Stderr: stderr.Bytes(),
			Err:    err,
		}
	}
	return &stdout, nil
}
