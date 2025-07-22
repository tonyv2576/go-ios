package run

import (
	"errors"
	"strings"
)

func FindCertificates() ([]*Certificate, error) {
	// get a list of every code signing identity available on the users device.
	out, err := runCommand("security", "find-identity", "-v", "-p", "codesigning")
	if err != nil {
		return nil, err
	}

	results := []*Certificate{}

	for _, line := range strings.Split(out, "\n") {
		cert := decodeCertificate(line)
		if cert != nil {
			results = append(results, cert)
		}

	}

	if len(results) <= 0 {
		return nil, errors.New("no code signing certificates found")
	}

	return results, nil
}

func ResignBundle(appName, certHash string) error {
	realError := func(msg string) bool {
		outMsg := strings.TrimSpace(msg)
		promptSuccess := "replacing existing signature"
		// output isnt empty on success so we have to check manually
		// note to self: if possible, find a better way to check for errors
		if len(outMsg) > 0 {
			if strings.HasSuffix(outMsg, promptSuccess) {
				return false
			}
		}
		return true
	}

	// codesign the app executable
	out, err := runCommand("codesign", "-f", "-s", certHash, "--entitlements", "entitlements.plist", appName+"/main")
	if err != nil {
		return err
	}
	if len(out) > 0 && realError(out) {

		return errors.New(out)
	}

	// codesign the whole bundle
	out, err = runCommand("codesign", "-f", "-s", certHash, "--entitlements", "entitlements.plist", appName)
	if err != nil {
		return err
	}
	if len(out) > 0 && realError(out) {
		return errors.New(out)
	}

	// both must be codesigned for the app to build.
	// note: codesigning the bundle alone will still pass the verification check
	// but the bundle still it won't be installable

	return nil
}

func VerifyBundle(appName string) (bool, error) {
	// verify the app was properly signed
	out, err := runCommand("codesign", "--verify", "--deep", "--strict", "--verbose=2", appName)
	if err != nil {
		return false, err
	}

	// note: very hacky way of checking if the command was a success. do better in the future
	validOnDisk, meetsRequirement := false, false
	promptDisk := "valid on disk"
	promptReq := "satisfies its designated requirement"

	for _, v := range strings.Split(out, "\n") {
		values := strings.Split(v, ":")

		if len(values) >= 2 {
			// make
			msg := strings.TrimSpace(values[1])
			if strings.EqualFold(msg, promptDisk) {
				validOnDisk = true
			} else if strings.EqualFold(msg, promptReq) {
				meetsRequirement = true
			} else {
				// any other output will be considered an error
				return false, errors.New(msg)
			}
		}
	}
	if !validOnDisk || !meetsRequirement {
		// failed without reason... hopefully we never run into this
		return false, errors.New("failed to verify bundle")
	}

	return true, nil
}
