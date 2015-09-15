package specsValidator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opencontainers/specs"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	// Path to config file inside the bundle
	ConfigFile  = "config.json"
	RuntimeFile = "runtime.json"
	// Path to rootfs directory inside the bundle
	RootfsDir = "rootfs"
)

type Bundle struct {
	Config  specs.Spec
	Runtime specs.RuntimeSpec

	Rootfs string
}

type LinuxBundle struct {
	Config  specs.LinuxSpec
	Runtime specs.LinuxRuntimeSpec

	Rootfs string
}

func ReadFile(file_url string) (content string, err error) {
	_, err = os.Stat(file_url)
	if err != nil {
		fmt.Println("cannot find the file ", file_url)
		return content, err
	}
	file, err := os.Open(file_url)
	defer file.Close()
	if err != nil {
		fmt.Println("cannot open the file ", file_url)
		return content, err
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	content = buf.String()

	return content, nil
}

func OSDetect(bundlePath string) string {
	content, err := ReadFile(path.Join(bundlePath, ConfigFile))
	if err == nil {
		var s specs.Spec
		err = json.Unmarshal([]byte(content), &s)
		if err == nil {
			return s.Platform.OS
		}
	}
	return ""
}

func FilesValid(bundlePath string, msgs []string) (bool, []string) {
	valid := true
	fi, err := os.Stat(bundlePath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Error accessing bundle: %v", err))
		return false, msgs
	} else {
		if !fi.IsDir() {
			msgs = append(msgs, fmt.Sprintf("Given path %s is not a directory", bundlePath))
			return false, msgs
		}
	}

	configPath := path.Join(bundlePath, ConfigFile)
	_, err = os.Stat(configPath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Error accessing %s: %v", ConfigFile, err))
		valid = false && valid
	}

	runtimePath := path.Join(bundlePath, RuntimeFile)
	_, err = os.Stat(runtimePath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Error accessing %s: %v", RuntimeFile, err))
		valid = false && valid
	}

	rootfsPath := path.Join(bundlePath, RootfsDir)
	fi, err = os.Stat(rootfsPath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Error accessing %s: %v", RootfsDir, err))
		valid = false && valid
	} else {
		if !fi.IsDir() {
			msgs = append(msgs, fmt.Sprintf("Given path %s is not a directory", rootfsPath))
			valid = false && valid
		}
	}

	return valid, msgs
}

func BundleValid(bundlePath string, msgs []string) (bool, []string) {
	var bundle Bundle
	valid, msgs := FilesValid(bundlePath, msgs)
	if valid == false {
		return valid, msgs
	}

	content, err := ReadFile(path.Join(bundlePath, ConfigFile))
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Cannot read %s", ConfigFile))
		valid = false && valid
	} else {
		err = json.Unmarshal([]byte(content), &(bundle.Config))
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Cannot parse %s", ConfigFile))
			valid = false && valid
		}
	}

	content, err = ReadFile(path.Join(bundlePath, RuntimeFile))
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Cannot read %s", RuntimeFile))
		valid = false && valid
	} else {
		err = json.Unmarshal([]byte(content), &(bundle.Runtime))
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Cannot parse %s", RuntimeFile))
			valid = false && valid
		}
	}

	bundle.Rootfs = path.Join(bundlePath, RootfsDir)
	if valid == false {
		return valid, msgs
	}

	ret, msgs := SpecValid(bundle.Config, bundle.Runtime, bundle.Rootfs, msgs)
	valid = ret && valid

	return valid, msgs
}

func LinuxBundleValid(bundlePath string, msgs []string) (bool, []string) {
	var bundle LinuxBundle
	valid, msgs := FilesValid(bundlePath, msgs)
	if valid == false {
		return valid, msgs
	}

	content, err := ReadFile(path.Join(bundlePath, ConfigFile))
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Cannot read %s", ConfigFile))
		valid = false && valid
	} else {
		err = json.Unmarshal([]byte(content), &(bundle.Config))
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Cannot parse %s", ConfigFile))
			valid = false && valid
		}
	}

	/*
		content, err = ReadFile(path.Join(bundlePath, RuntimeFile))
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Cannot read %s", RuntimeFile))
			valid = false && valid
		} else {
			err = json.Unmarshal([]byte(content), &(bundle.Runtime))
			if err != nil {
				msgs = append(msgs, fmt.Sprintf("Cannot parse %s", RuntimeFile))
				valid = false && valid
			}
		}
	*/
	bundle.Rootfs = path.Join(bundlePath, RootfsDir)
	if valid == false {
		return valid, msgs
	}

	ret, msgs := LinuxSpecValid(bundle.Config, bundle.Runtime, bundle.Rootfs, msgs)
	valid = ret && valid
	return valid, msgs
}

func StringValid(key string, content string, msgs []string) (bool, []string) {
	valid := true
	if len(content) == 0 {
		valid = false
		msg := fmt.Sprintf("%s is missing", key)
		msgs = append(msgs, msg)
	}
	return valid, msgs
}

func checkSemVer(version string, msgs []string) (bool, []string) {
	valid, msgs := StringValid("Spec.Version", version, msgs)
	if valid == false {
		return valid, msgs
	}
	if version == specs.Version {
		return true, msgs
	}

	str := strings.Split(version, ".")
	if len(str) != 3 {
		valid = false
	} else {
		for index := 0; index < len(str); index++ {
			i, err := strconv.Atoi(str[index])
			if err != nil {
				valid = false
				break
			} else {
				if i < 0 {
					valid = false
					break
				}
			}
		}
	}
	if valid == false {
		msg := fmt.Sprintf("%s is not a valid version format, please read 'SemVer v2.0.0'", version)
		msgs = append(msgs, msg)
	}
	return valid, msgs
}
