package main

import "errors"

var (
	ErrNoArg       error = errors.New("must choose one of the following arguments: mode-build, mode-build-std, mode-sign, mode-install")
	ErrFlagBundle  error = errors.New("bundle flag is required")
	ErrFlagProfile error = errors.New("profile flag is required")
	ErrFlagCert    error = errors.New("cert flag is required on devices with multiple certificates")

	ErrUnverified error = errors.New("failed to verify build")

	ErrNoBuild        error = errors.New("no project build found")
	ErrNoCertificates error = errors.New("no code signing certificates found")
)
