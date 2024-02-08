package handlers

import (
	"fmt"
	"strings"
)

type CmdError struct {
	Cmd    []string
	Stderr []byte
	Err    error
}

func (err *CmdError) Error() string {
	return fmt.Sprintf("Command '%s' failed: %v\nStderr: %v", strings.Join(err.Cmd, " "), err.Err, string(err.Stderr))
}
