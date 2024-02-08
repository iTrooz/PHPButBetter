package main

import (
	"bytes"
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

	gpp_cmd := exec.Command("javac", filepath, "-d", tmpFolder)
	var gppStderr bytes.Buffer
	gpp_cmd.Stderr = &gppStderr
	err = gpp_cmd.Run()
	if err != nil {
		log.Errorf("Command stderr (corresponding error below):\n%v", gppStderr.String())
		return fmt.Errorf("Failed to run javac command: %w", err)
	}

	cmd := exec.Command("java", "-cp", tmpFolder, className)
	var cmdStdout, cmdStderr bytes.Buffer
	cmd.Stdout = &cmdStdout
	cmd.Stderr = &cmdStderr
	err = cmd.Run()
	if err != nil {
		log.Errorf("Command stderr (corresponding error below):\n%v", cmdStderr.String())
		return fmt.Errorf("Failed to run java code: %w", err)
	}

	stdoutSize := cmdStdout.Len()
	writtenBytes, err := w.Write(cmdStdout.Bytes())
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
