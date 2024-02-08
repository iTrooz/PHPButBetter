package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
)

func CHandler(w http.ResponseWriter, filepath string) error {
	return CCppHandler("gcc", w, filepath)
}
func CppHandler(w http.ResponseWriter, filepath string) error {
	return CCppHandler("g++", w, filepath)
}

func CCppHandler(compiler string, w http.ResponseWriter, filepath string) error {

	tmpFolder, err := os.MkdirTemp("", "phpbutbetter")
	if err != nil {
		return fmt.Errorf("Failed to create temporary folder: %w", err)
	}

	compiledCodePath := path.Join(tmpFolder, "a.out")
	_, err = RunCmd(exec.Command(compiler, filepath, "-o", compiledCodePath))
	if err != nil {
		return err
	}

	stdout, err := RunCmd(exec.Command(compiledCodePath))
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