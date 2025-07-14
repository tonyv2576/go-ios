package run

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

func runCommand(name string, params ...string) (string, error) {
	localDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return runCommandEx(name, localDir, nil, params...)
}

func runCommandOut(out io.Writer, name string, params ...string) error {
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
	c.Stderr = output
	if out != nil {
		c.Stdout = out
	} else {
		c.Stdout = output
	}

	if err := c.Start(); err != nil {
		return "", err
	}
	if err := c.Wait(); err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}
