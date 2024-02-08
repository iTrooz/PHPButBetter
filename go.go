package main

import (
	"fmt"
	"net/http"
)

func GoHandler(w http.ResponseWriter, filepath string) error {
	stdout, err := RunCmd("go", "run", filepath)
	if err != nil {
		return err
	}

	stdoutSize := stdout.Len()
	writtenBytes, err := w.Write(stdout.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to write compiled code stdout as response: %w", err)
	}
	if writtenBytes != stdoutSize {
		return fmt.Errorf("Failed to fully write compiled code stdout as response: wrote %v/%v", writtenBytes, stdoutSize)
	}

	return nil
}
