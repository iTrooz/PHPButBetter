package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
)

func CMakeHandler(w http.ResponseWriter, filepath string) error {
	cmd := exec.Command("cmake", "-P", filepath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return &CmdError{
			Cmd:    cmd.Args,
			Stderr: stderr.Bytes(),
			Err:    err,
		}
	}

	// CMake prints message() calls to stderr
	stderrSize := stderr.Len()
	writtenBytes, err := w.Write(stderr.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to write code output as response: %w", err)
	}
	if writtenBytes != stderrSize {
		return fmt.Errorf("Failed to fully write code output as response: wrote %v/%v", writtenBytes, stderrSize)
	}

	return nil
}
