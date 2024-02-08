package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
)

func GoHandler(w http.ResponseWriter, filepath string) error {

	cmd := exec.Command("go", "run", filepath)
	var cmdStdout bytes.Buffer
	cmd.Stdout = &cmdStdout
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to run compiled code: %w", err)
	}

	stdoutSize := cmdStdout.Len()
	writtenBytes, err := w.Write(cmdStdout.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to write compiled code stdout as response: %w", err)
	}
	if writtenBytes != stdoutSize {
		return fmt.Errorf("Failed to fully write compiled code stdout as response: wrote %v/%v", writtenBytes, stdoutSize)
	}

	return nil
}
