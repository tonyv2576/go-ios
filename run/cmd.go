package run

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

func runCommand(name string, params ...string) (string, error) {
	// short hand for running commands in the current directory
	localDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return runCommandEx(name, localDir, nil, params...)
}

func runCommandOut(out io.Writer, name string, params ...string) error {
	// shorthand for running commands with an io.Writer instead of a string output
	localDir, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = runCommandEx(name, localDir, out, params...)
	return err
}

func runCommandEx(name, dir string, out io.Writer, params ...string) (string, error) {
	c := exec.Command(name, params...)
	if len(dir) > 0 {
		c.Dir = dir
	}

	output := bytes.NewBuffer(nil)
	errOutput := bytes.NewBuffer(nil)

	c.Stdout = output
	c.Stderr = errOutput
	// use the writer if one was passed in
	if out != nil {
		c.Stdout = out
	}

	// return the stderr output on fail
	if err := c.Start(); err != nil {
		return output.String(), errors.New(errOutput.String())
	}
	if err := c.Wait(); err != nil {
		return output.String(), errors.New(errOutput.String())
	}

	// get rid of unnecessary white space
	result := strings.TrimSpace(output.String())
	if len(result) <= 0 {
		// fallback because some commands use stderr only
		result = strings.TrimSpace(errOutput.String())
	}
	return result, nil
}
