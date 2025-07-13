package run

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func runCommand(name string, params ...string) (string, error) {
	localDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return runCommandAt(name, localDir, params...)
}

func runCommandAt(name, dir string, params ...string) (string, error) {
	c := exec.Command(name, params...)
	if len(dir) > 0 {
		c.Dir = dir
	}

	output := bytes.NewBuffer(nil)
	c.Stdout = output

	if err := c.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}
