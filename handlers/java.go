package handlers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

func JavaHandler(w http.ResponseWriter, filepath string) error {

	tmpFolder, err := os.MkdirTemp("", "phpbutbetter")
	if err != nil {
		return fmt.Errorf("Failed to create temporary folder: %w", err)
	}

	className := strings.TrimSuffix(path.Base(filepath), ".java")

	_, err = runCmd(exec.Command("javac", filepath, "-d", tmpFolder))
	if err != nil {
		return err
	}

	stdout, err := runCmd(exec.Command("java", "-cp", tmpFolder, className))
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

	err = os.RemoveAll(tmpFolder)
	if err != nil {
		return fmt.Errorf("Failed to remove temporary folder: %w", err)
	}

	return nil
}
