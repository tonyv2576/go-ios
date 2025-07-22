package run

import (
	"errors"
	"strings"
)

func IsInstalled(cmdname string) (bool, error) {
	// check if a command exists
	out, err := runCommand("type", cmdname)
	if err != nil {
		return false, err
	}
	out = strings.ToLower(out)
	cmdname = strings.ToLower(cmdname)

	return strings.HasPrefix(out, cmdname), nil
}

// i considered just installing ios-deploy myself if it wasnt found but that would still require the user to install brew

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
	// check for dependencies before running any commands
	// note: most of these should be installed onto your mac by default
	check := func(cmd string) error {
		found, err := IsInstalled(cmd)
		if err != nil {
			return err
		} else if !found {
			// return a pretty error to the user
			return errors.New("unable to find the following command: " + cmd)
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
