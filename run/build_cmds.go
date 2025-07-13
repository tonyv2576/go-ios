package run

import (
	"errors"
	"os"
	"path"
	"strings"
)

func BundleInfo() (*Bundle, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	b := Bundle{}
	b.ProjectPath = dir
	b.Bundle.Name = path.Base(dir) + ".app"
	b.Bundle.Path = path.Join(dir, b.Bundle.Name)

	return &b, nil
}

func BuildBundle(bundleId string, codePath string) error {
	out, err := runCommand("gomobile", "build", "-target=ios", "-bundleid="+bundleId, codePath)
	if err != nil {
		return err
	}

	if len(out) > 0 {
		return errors.New(out)
	}

	return nil
}

func ConvertPlist(bundlePath string) error {
	plistPath := path.Join(bundlePath, "Info.plist")
	out, err := runCommand("plutil", "-convert", "xml1", plistPath)
	if err != nil {
		return err
	}

	if len(out) > 0 {
		return errors.New(out)
	}

	return nil
}

func EmbedMobileProfile(bundleName, profilePath string) error {
	out, err := runCommand("cp", profilePath, bundleName+"/embedded.mobileprovision")
	if err != nil {
		return err
	}
	if len(out) > 0 {
		return errors.New(out)
	}

	return nil
}
func ExtractEntitlements(profilePath string) error {
	out, err := runCommand("security", "cms", "-D", "-i", profilePath, ">", "temp.plist")
	if err != nil {
		return err
	}
	if len(out) > 0 {
		return errors.New(out)
	}

	out, err = runCommand("/usr/libexec/PlistBuddy", "-x", "-c", "'Print:Entitlements'", "temp.plist", ">", "entitlements.plist")
	if err != nil {
		return err
	}
	if len(out) > 0 {
		return errors.New(out)
	}

	return nil
}

func RemoveTempFiles(projectPath string) {
	os.Remove(path.Join(projectPath, "entitlements.plist"))
	os.Remove(path.Join(projectPath, "temp.plist"))
}

func DeployBundle(bundleName string) (*BuildDetails, error) {
	out, err := runCommand("ios-deploy", "--bundle", bundleName)
	if err != nil {
		return nil, err
	}

	phase := -1
	progress := "---"
	action := ""
	errMsg := ""

	for _, v := range strings.Split(out, "\n") {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "------") {
			if strings.Contains(strings.ToLower(v), "install phase") {
				phase = 0
			}
		}
		if strings.HasPrefix(v, "[") {
			if phase >= 0 && len(v) >= 6 {
				progress = strings.TrimSpace(v[:6])
				progress = strings.TrimPrefix(strings.TrimSuffix(progress, "]"), "[")

				action = strings.TrimSpace(strings.TrimPrefix(v, progress))
			}
		} else {
			errMsg = v
			break
		}
	}

	if len(errMsg) > 0 {
		err = errors.New(errMsg)
		details := &BuildDetails{
			Progress: progress,
			Action:   action,
			Error:    err,
		}

		return details, errors.New(details.String())
	}

	return &BuildDetails{
		Progress: progress,
		Action:   "",
		Error:    nil,
	}, nil
}
