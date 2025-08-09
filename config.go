package main

import (
	"fmt"
	"go-ios/run"
)

type Config struct {
	Bundle       string
	Certificate  string
	ProfilePath  string
	CodePath     string
	AppendedPath string

	Simulator bool
	Unsigned  bool
	NoInstall bool
	Silent    bool

	Mode int // 0 std-build, build, 1 sign, 2 install
}

func (c *Config) Do() error {
	if err := run.CheckDependencies(); err != nil {
		return err
	}

	info, err := run.BundleInfo()
	if err != nil {
		return err
	}

	switch c.Mode {
	case 1:
		// export
		if len(c.Bundle) <= 0 {
			return ErrFlagBundle
		}
		if len(c.ProfilePath) > 0 {
			c.Unsigned = true
		}

		if err := c.cmdBuild(info); err != nil {
			return err
		}
		if c.Unsigned {
			if err := c.cmdCodesign(info); err != nil {
				return err
			}
		}
		if c.NoInstall {
			c.log("App successfully exported!")
		} else {
			c.log("Deploying to device...")
			if _, err := run.DeployBundle(info.Bundle.Name); err != nil {
				return err
			}
			c.log("App successfully installed!")
		}
	case 2:
		// build
		// once we edit the build, it will need to be co-signed again.
		if len(c.Bundle) <= 0 {
			return ErrFlagBundle
		}

		if err := c.cmdBuild(info); err != nil {
			return err
		}
		c.log("Build successful!")

	case 3:
		// codesign
		if !info.Initialized() {
			return ErrNoBuild
		}
		if len(c.ProfilePath) <= 0 {
			return ErrFlagProfile
		}

		if err := c.cmdCodesign(info); err != nil {
			return err
		}
		c.log("Successfully signed!")

	case 4:
		// install/deploy
		if !info.Initialized() {
			return ErrNoBuild
		}
		if _, err := run.DeployBundle(info.Bundle.Name); err != nil {
			return err
		}
		c.log("App successfully installed!")
	default:
		return ErrNoArg
	}

	return nil
}

func (c *Config) cmdBuild(info *run.Bundle) error {
	c.log("Building bundle...")
	if err := run.BuildBundle(c.Bundle, c.CodePath, c.Simulator); err != nil {
		return err
	}

	// unsigned build
	if c.Unsigned {
		if len(c.AppendedPath) > 0 {
			c.log("Updating Info.Plist with additional keys...")
			if err := run.AppendToPlist(info.Bundle.Path, conf.AppendedPath); err != nil {
				return err
			}
		} else {
			c.log("Re-encoding Info.Plist as XML...")
			if err := run.ConvertPlist(info.Bundle.Path); err != nil {
				return err
			}
		}
	}
	// elsewise, gomobile signs them using xcode's CLI

	return nil
}

func (c *Config) cmdCodesign(info *run.Bundle) error {
	hash := c.Certificate

	// if no hash is provided
	if len(hash) <= 0 {
		// use the findcertificates function and use the one in there.
		// note: if more than one exists, one must be provided manually anyway
		c.log("Looking for signing certificates...")
		if certs, err := run.FindCertificates(); err != nil {
			return err
		} else if len(certs) <= 0 {
			return ErrNoCertificates
		} else if len(certs) > 1 {
			return ErrFlagCert
		} else {
			hash = certs[0].Hash
		}
	}
	// rid ourselves of those pesky temporary files (even if the function errors)
	defer run.RemoveTempFiles(info.ProjectPath)

	c.log("Embedding mobile provisioning profile...")
	if err := run.EmbedMobileProfile(info.Bundle.Name, c.ProfilePath); err != nil {
		return err
	}

	c.log("Extracting profile entitlements...")
	if err := run.ExtractEntitlements(c.ProfilePath); err != nil {
		return err
	}

	c.log("Re-signing bundle...")
	if err := run.ResignBundle(info.Bundle.Name, hash); err != nil {
		return err
	}

	c.log("Verifying bundle...")
	if suceess, err := run.VerifyBundle(info.Bundle.Name); err != nil {
		return err
	} else if !suceess {
		return ErrUnverified
	}

	return nil
}

func (c *Config) log(message string, args ...any) {
	if !c.Silent {
		fmt.Printf(message+"\n", args...)
	}
}
