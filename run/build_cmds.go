package run

import (
	"errors"
	"os"
	"path"
	"strings"

	"howett.net/plist"
)

func BundleInfo() (*Bundle, error) {
	// get the full path of the bundle
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
	// note: simulator support is buggy and does not currently work with cgo
	target := "ios"
	if simulator {
		target += "simulator"
	}
	// this will not work unless your project imports "gomobile" directly and
	// importing it twice causes linker issues so if your library already has it's
	// own build tool, use that instead. (Ex. fyne)
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

	// locate Info.plist
	plistPath := bundlePath + "/Info.plist"
	plistFile, err := os.Open(plistPath)
	if err != nil {
		return err
	}
	defer plistFile.Close()

	// decode and re-encode as xml1
	decoder := plist.NewDecoder(plistFile)
	if err := decoder.Decode(&plistData); err != nil {
		return err
	}

	f, err := os.OpenFile(plistPath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	// overwrite file
	encoder := plist.NewEncoder(f)
	if err := encoder.Encode(plistData); err != nil {
		return err
	}

	return nil
}

func AppendToPlist(bundlePath, appendedPath string) error {
	var plistData map[string]any
	var appendedData map[string]any

	// locate plist
	plistPath := bundlePath + "/Info.plist"
	plistFile, err := os.Open(plistPath)
	if err != nil {
		return err
	}
	defer plistFile.Close()

	// locate user provided plist file
	// plist preferred but xml works as-well
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

	// append (and replace) keys and values to the app's plist
	for key, value := range appendedData {
		plistData[key] = value
	}

	f, err := os.OpenFile(plistPath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	// overwrite file and re-encode as binary. xml format would've worked too but you should need to
	// edit the file manually if you're using an external plist
	encoder := plist.NewBinaryEncoder(f)
	if err := encoder.Encode(plistData); err != nil {
		return err
	}

	return nil
}

func EmbedMobileProfile(bundleName, profilePath string) error {
	// copy the data from the mobile provision into our bundle's embedded provision
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
	// create a temporary file
	tempOut, err := os.Create("temp.plist")
	if err != nil {
		return err
	}
	defer tempOut.Close()

	// extract plist from the provided mobileprovision file
	err = runCommandOut(tempOut, "security", "cms", "-D", "-i", profilePath)
	if err != nil {
		return err
	}

	// another temporary file
	entOut, err := os.Create("entitlements.plist")
	if err != nil {
		return err
	}
	defer entOut.Close()

	// extract the entitlements from the temporary plist file we made
	err = runCommandOut(entOut, "/usr/libexec/PlistBuddy", "-x", "-c", "Print:Entitlements", "temp.plist")
	if err != nil {
		return err
	}

	return nil
}

func RemoveTempFiles(projectPath string) {
	// clean up those pesky temporary files
	os.Remove(path.Join(projectPath, "entitlements.plist"))
	os.Remove(path.Join(projectPath, "temp.plist"))
}

func DeployBundle(bundleName string) (*BuildDetails, error) {
	// deploy with ios-deploy. not necessary, you can install it manually if you so choose but currently
	// this tool requires it still be installed.
	// note to future self: give users the ability to use the tool without ios-deploy
	out, err := runCommand("ios-deploy", "--bundle", bundleName)
	if err != nil {
		return nil, err
	}

	progress := "---"
	action := ""
	errMsg := ""

	// note: find a better way to log errors
	// this "works" but i doubt it's fool proofed
	for _, v := range strings.Split(out, "\n") {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "-") {
			continue
		}
		if strings.HasPrefix(v, "[") {
			// keep track of the completion state
			if len(v) >= 6 {
				progress = strings.TrimSpace(v[:6])
				progress = strings.TrimPrefix(strings.TrimSuffix(progress, "]"), "[")

				action = strings.TrimSpace(strings.TrimPrefix(v, progress))
			}
		} else {
			// return error
			errMsg = v
			break
		}
	}

	// return error
	if len(errMsg) > 0 {
		err = errors.New(errMsg)
		details := &BuildDetails{
			Progress: progress,
			Action:   action,
			Error:    err,
		}

		return details, errors.New(details.String())
	}

	// return nil
	return &BuildDetails{
		Progress: progress,
		Action:   "",
		Error:    nil,
	}, nil
}
