package main

import (
	"go-ios/run"
)

type Config struct {
	Bundle      string
	Certificate string
	ProfilePath string
	CodePath    string

	Mode int // 0 std-build, build, 1 sign, 2 install
}

func (c *Config) Do() error {
	info, err := run.BundleInfo()
	if err != nil {
		return err
	}

	switch c.Mode {
	case 1:
		// standard build
		if len(c.Bundle) <= 0 {
			return ErrFlagBundle
		}
		if err := run.BuildBundle(c.Bundle, info.ProjectPath); err != nil {
			return err
		}

	case 2:
		// editable build
		if len(c.Bundle) <= 0 {
			return ErrFlagBundle
		}

		if err := run.BuildBundle(c.Bundle, c.CodePath); err != nil {
			return err
		}
		if err := run.ConvertPlist(info.Bundle.Path); err != nil {
			return err
		}
	case 3:
		// codesign build
		if !info.Initialized() {
			return ErrNoBuild
		}
		if len(c.ProfilePath) <= 0 {
			return ErrFlagProfile
		}

		hash := c.Certificate

		if len(hash) <= 0 {
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
		defer run.RemoveTempFiles(info.ProjectPath)

		if err := run.EmbedMobileProfile(info.Bundle.Name, c.ProfilePath); err != nil {
			return err
		}
		if err := run.ExtractEntitlements(c.ProfilePath); err != nil {
			return err
		}
		if err := run.ResignBundle(info.Bundle.Name, hash); err != nil {
			return err
		}
		if suceess, err := run.VerifyBundle(info.Bundle.Name); err != nil {
			return err
		} else if !suceess {
			return ErrUnverified
		}
	case 4:
		// install/deploy build
		if !info.Initialized() {
			return ErrNoBuild
		}
		if _, err := run.DeployBundle(info.Bundle.Name); err != nil {
			return err
		}

	default:
		return ErrNoArg
	}

	return nil
}
