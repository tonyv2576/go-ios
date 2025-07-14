package run

import (
	"errors"
	"strings"
)

func FindCertificates() ([]*Certificate, error) {
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

		if len(outMsg) > 0 {
			if strings.HasSuffix(outMsg, promptSuccess) {
				return false
			}
		}
		return true
	}
	out, err := runCommand("codesign", "-f", "-s", certHash, "--entitlements", "entitlements.plist", appName+"/main")
	if err != nil {
		return err
	}
	if len(out) > 0 && realError(out) {

		return errors.New(out)
	}

	out, err = runCommand("codesign", "-f", "-s", certHash, "--entitlements", "entitlements.plist", appName)
	if err != nil {
		return err
	}
	if len(out) > 0 && realError(out) {
		return errors.New(out)
	}

	return nil
}

func VerifyBundle(appName string) (bool, error) {
	out, err := runCommand("codesign", "--verify", "--deep", "--strict", "--verbose=2", appName)
	if err != nil {
		return false, err
	}

	validOnDisk, meetsRequirement := false, false
	promptDisk := "valid on disk"
	promptReq := "satisfies its designated requirement"

	for _, v := range strings.Split(out, "\n") {
		values := strings.Split(v, ":")

		if len(values) >= 2 {
			msg := strings.TrimSpace(values[1])
			if strings.EqualFold(msg, promptDisk) {
				validOnDisk = true
			} else if strings.EqualFold(msg, promptReq) {
				meetsRequirement = true
			} else {
				return false, errors.New(msg)
			}
		}
	}
	if !validOnDisk || !meetsRequirement {
		return false, errors.New("failed to verify bundle")
	}

	return true, nil
}
