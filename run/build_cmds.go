package run

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"howett.net/plist"
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

func BuildBundle(bundleId string, codePath string, simulator bool) error {
	target := "ios"
	if simulator {
		target += "simulator"
	}
	out, err := runCommand("gomobile", "build", "-target="+target, "-bundleid="+bundleId, codePath)
	if err != nil {
		return err
	}

	if len(out) > 0 {
		return errors.New(out)
	}

	return nil
}

func ConvertPlist(bundlePath string) error {
	var plistData map[string]any

	plistPath := bundlePath + "/Info.plist"
	plistFile, err := os.Open(plistPath)
	if err != nil {
		return err
	}
	defer plistFile.Close()
	decoder := plist.NewDecoder(plistFile)
	if err := decoder.Decode(&plistData); err != nil {
		return err
	}

	f, err := os.OpenFile(plistPath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := plist.NewEncoder(f)
	if err := encoder.Encode(plistData); err != nil {
		return err
	}

	return nil
}

func AppendToPlist(bundlePath, appendedPath string) error {
	var plistData map[string]any
	var appendedData map[string]any

	plistPath := bundlePath + "/Info.plist"
	plistFile, err := os.Open(plistPath)
	if err != nil {
		return err
	}
	defer plistFile.Close()
	appendedFile, err := os.Open(appendedPath)
	if err != nil {
		return err
	}
	defer appendedFile.Close()

	decoder := plist.NewDecoder(plistFile)
	if err := decoder.Decode(&plistData); err != nil {
		return err
	}

	decoder = plist.NewDecoder(appendedFile)
	if err := decoder.Decode(&appendedData); err != nil {
		return err
	}

	for key, value := range appendedData {
		plistData[key] = value
	}

	f, err := os.OpenFile(plistPath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := plist.NewBinaryEncoder(f)
	if err := encoder.Encode(plistData); err != nil {
		return err
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
	tempOut, err := os.Create("temp.plist")
	if err != nil {
		return err
	}
	defer tempOut.Close()

	err = runCommandOut(tempOut, "security", "cms", "-D", "-i", profilePath)
	if err != nil {
		return err
	}
	entOut, err := os.Create("entitlements.plist")
	if err != nil {
		return err
	}
	defer entOut.Close()

	err = runCommandOut(entOut, "/usr/libexec/PlistBuddy", "-x", "-c", "Print:Entitlements", "temp.plist")
	if err != nil {
		return err
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

	progress := "---"
	action := ""
	errMsg := ""

	for _, v := range strings.Split(out, "\n") {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "-") {
			continue
		}
		if strings.HasPrefix(v, "[") {
			if len(v) >= 6 {
				progress = strings.TrimSpace(v[:6])
				progress = strings.TrimPrefix(strings.TrimSuffix(progress, "]"), "[")

				action = strings.TrimSpace(strings.TrimPrefix(v, progress))
			}
		} else {
			fmt.Println(progress, action)
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
