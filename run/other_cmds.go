package run

import (
	"errors"
	"strings"
)

func IsInstalled(cmdname string) (bool, error) {
	out, err := runCommand("type", cmdname)
	if err != nil {
		return false, err
	}
	out = strings.ToLower(out)
	cmdname = strings.ToLower(cmdname)

	return strings.HasPrefix(out, cmdname), nil
}

// func InstallIosDeploy() error {
// 	found, err := IsInstalled("ios-deploy")
// 	if err != nil {
// 		return err
// 	}
// 	if !found {
// 		var stderr bytes.Buffer

// 		cmd := exec.Command("brew", "install", "ios-deploy", "--quiet", "--no-quarantine", "--formula")
// 		cmd.Env = append(os.Environ(), "HOMEBREW_NO_AUTO_UPDATE=1")
// 		cmd.Stderr = &stderr

// 		if err := cmd.Run(); err != nil {
// 			return errors.New(stderr.String())
// 		}
// 	}

// 	return nil
// }

func CheckDependencies() error {
	check := func(cmd string) error {
		found, err := IsInstalled(cmd)
		if err != nil {
			return err
		} else if !found {
			return errors.New("unable to find command: " + cmd)
		}
		return nil
	}

	if err := check("plutil"); err != nil {
		return err
	}
	if err := check("/usr/libexec/PlistBuddy"); err != nil {
		return err
	}
	if err := check("ios-deploy"); err != nil {
		return err
	}

	return nil
}
