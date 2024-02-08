package handlers

import (
	"bytes"
	"os/exec"
)

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
