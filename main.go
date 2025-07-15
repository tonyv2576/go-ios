package main

import (
	"flag"
	"fmt"
	"os"
)

var conf = &Config{}

func init() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'build', 'sign', or 'install' subcommands")
		os.Exit(1)
	}

	var bundle, profile, cert, appended string
	var unsigned, simulator bool

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&bundle, "bundle", "", "Bundle identifier")
	buildCmd.StringVar(&appended, "append", "", "Creates an unsigned build and with custom keys")
	buildCmd.BoolVar(&unsigned, "unsigned", false, "Creates an unsigned build with an editable plist file")
	buildCmd.BoolVar(&simulator, "simulator", false, "Creates a build compatible with ios simulators")

	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signCmd.StringVar(&profile, "profile", "", "Mobile provisioning profile")
	signCmd.StringVar(&cert, "cert", "", "Certificate to use")

	installCmd := flag.NewFlagSet("install", flag.ExitOnError)

	switch os.Args[1] {
	case "build":
		conf.Mode = 1
		buildCmd.Parse(os.Args[2:])

		if unsigned || len(appended) > 0 {
			conf.Mode = 2
		}
	case "sign":
		conf.Mode = 3
		signCmd.Parse(os.Args[2:])
	case "install":
		conf.Mode = 4
		installCmd.Parse(os.Args[2:])
	default:
		fmt.Println("unknown subcommand:", os.Args[1])
		os.Exit(1)
	}

	var targetDir string

	// capture trailing non-flag argument as targetDir
	args := os.Args[2:]
	for i, arg := range args {
		if arg == "--" {
			if i+1 < len(args) {
				targetDir = args[i+1]
			}
			break
		} else if len(arg) > 0 && arg[0] != '-' {
			// if it looks like a relative path and isn't parsed by flags
			if _, err := os.Stat(arg); err == nil || os.IsNotExist(err) {
				targetDir = arg
				break
			}
		}
	}

	conf.Bundle = bundle
	conf.ProfilePath = profile
	conf.Certificate = cert
	conf.CodePath = targetDir
	conf.AppendedPath = appended
	conf.Simulator = simulator
}

func main() {
	if err := conf.Do(); err != nil {
		fmt.Println(err)
	}
}
