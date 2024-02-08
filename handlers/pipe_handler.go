package handlers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func PipeWrapper(cmdName string) func(w http.ResponseWriter, filepath string) error {
	return func(w http.ResponseWriter, filepath string) error {
		return pipeHandler(cmdName, w, filepath)
	}
}

func pipeHandler(cmdName string, w http.ResponseWriter, filepath string) error {

	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}

	cmd := exec.Command(cmdName)
	cmd.Stdin = f
	stdout, err := runCmd(cmd)
	if err != nil {
		return err
	}

	stdoutSize := stdout.Len()
	writtenBytes, err := w.Write(stdout.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to write code stdout as response: %w", err)
	}
	if writtenBytes != stdoutSize {
		return fmt.Errorf("Failed to fully write code stdout as response: wrote %v/%v", writtenBytes, stdoutSize)
	}

	return nil
}
