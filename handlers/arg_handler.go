package handlers

import (
	"fmt"
	"net/http"
	"os/exec"
)

func ArgWrapper(cmdName string) func(w http.ResponseWriter, filepath string) error {
	return func(w http.ResponseWriter, filepath string) error {
		return ArgHandler(cmdName, w, filepath)
	}
}

func ArgHandler(cmdName string, w http.ResponseWriter, filepath string) error {
	stdout, err := runCmd(exec.Command(cmdName, filepath))
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
