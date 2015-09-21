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
	Rootfs  string
}

type LinuxBundle struct {
	Config  specs.LinuxSpec
	Runtime specs.LinuxRuntimeSpec
	Rootfs  string
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

func OSDetect(inputPath string) string {
	var configURL string
	fi, err := os.Stat(inputPath)
	if err == nil {
		if fi.IsDir() {
			configURL = path.Join(inputPath, ConfigFile)
		} else {
			configURL = inputPath
		}
	} else {
		return ""
	}
	content, err := ReadFile(configURL)
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

func OCTConfigValid(configPath string, msgs []string) (bool, []string) {
	valid := true
	os := OSDetect(configPath)
	if len(os) == 0 {
		msgs = append(msgs, "Cannot detect OS in the config.json under the bundle, or maybe miss `config.json`.")
		return false, msgs
	}
	if os == "linux" {
		var ls specs.LinuxSpec
		var rt specs.LinuxRuntimeSpec
		content, _ := ReadFile(configPath)
		json.Unmarshal([]byte(content), &ls)
		valid, msgs = LinuxSpecValid(ls, rt, "", msgs)
	} else {
		var s specs.Spec
		var rt specs.RuntimeSpec
		content, _ := ReadFile(configPath)
		json.Unmarshal([]byte(content), &s)
		valid, msgs = SpecValid(s, rt, "", msgs)
	}
	return valid, msgs
}

func OCTRuntimeValid(runtimePath string, os string, msgs []string) (bool, []string) {
	valid := true
	content, err := ReadFile(runtimePath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("Cannot read %s", runtimePath))
		return false, msgs
	}
	if os == "linux" {
		var lrt specs.LinuxRuntimeSpec
		err = json.Unmarshal([]byte(content), &lrt)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Cannot parse %s", runtimePath))
			valid = false
		} else {
			valid, msgs = LinuxRuntimeSpecValid(lrt, "", msgs)
		}
	} else {
		var rt specs.RuntimeSpec
		err = json.Unmarshal([]byte(content), &rt)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("Cannot parse %s", runtimePath))
			valid = false
		} else {
			valid, msgs = RuntimeSpecValid(rt, "", msgs)
		}
	}
	return valid, msgs
}

func OCTBundleValid(bundlePath string, msgs []string) (bool, []string) {
	valid, msgs := FilesValid(bundlePath, msgs)
	if valid == false {
		return valid, msgs
	}

	os := OSDetect(bundlePath)
	if len(os) == 0 {
		msgs = append(msgs, "Cannot detect OS in the config.json under the bundle, or maybe miss `config.json`.")
		return false, msgs
	}

	if os == "linux" {
		valid, msgs = LinuxBundleValid(bundlePath, msgs)
	} else {
		valid, msgs = BundleValid(bundlePath, msgs)
	}

	return valid, msgs
}

func BundleValid(bundlePath string, msgs []string) (bool, []string) {
	var bundle Bundle
	valid := true

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

	if valid == false {
		return valid, msgs
	}

	bundle.Rootfs = path.Join(bundlePath, RootfsDir)

	ret, msgs := SpecValid(bundle.Config, bundle.Runtime, bundle.Rootfs, msgs)
	valid = ret && valid

	ret, msgs = RuntimeSpecValid(bundle.Runtime, bundle.Rootfs, msgs)
	valid = ret && valid

	return valid, msgs
}

func LinuxBundleValid(bundlePath string, msgs []string) (bool, []string) {
	var bundle LinuxBundle
	valid := true

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

	if valid == false {
		return valid, msgs
	}

	bundle.Rootfs = path.Join(bundlePath, RootfsDir)

	ret, msgs := LinuxSpecValid(bundle.Config, bundle.Runtime, bundle.Rootfs, msgs)
	valid = ret && valid

	ret, msgs = LinuxRuntimeSpecValid(bundle.Runtime, bundle.Rootfs, msgs)
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
