package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
)

func CppHandler(w http.ResponseWriter, filepath string) error {

	tmpFolder, err := os.MkdirTemp("", "phpbutbetter")
	if err != nil {
		return fmt.Errorf("Failed to create temporary folder: %w", err)
	}

	compiledCodePath := path.Join(tmpFolder, "a.out")
	gpp_cmd := exec.Command("g++", filepath, "-o", compiledCodePath)
	var gppStderr bytes.Buffer
	gpp_cmd.Stderr = &gppStderr
	err = gpp_cmd.Run()
	if err != nil {
		log.Errorf("Command stderr (corresponding error below):\n%v", gppStderr.String())
		return fmt.Errorf("Failed to run g++ command: %w", err)
	}

	cmd := exec.Command(compiledCodePath)
	var cmdStdout bytes.Buffer
	cmd.Stdout = &cmdStdout
	err = cmd.Run()
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

	err = os.RemoveAll(tmpFolder)
	if err != nil {
		return fmt.Errorf("Failed to remove temporary folder: %w", err)
	}

	return nil
}
