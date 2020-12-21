// Package sysutils contains system-related utilities.
package sysutils

import (
	"os/exec"

	"gopkg.in/alessio/shellescape.v1"
)

func escapeArgs(args []string) []string {
	output := make([]string, len(args))

	for i, arg := range args {
		output[i] = shellescape.Quote(arg)
	}

	return output
}

// RunCommandInDir runs the command in the directory pointed to by dir,
// after escaping the arguments.
// This function returns stdout and err if any.
func RunCommandInDir(bin string, args []string, dir string) ([]byte, error) {
	escapedArgs := escapeArgs(args)
	cmd := exec.Command(bin, escapedArgs...)
	cmd.Dir = dir
	return cmd.Output()
}
