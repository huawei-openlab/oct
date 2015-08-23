// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	configconvert "./../../source/configconvert"
	runcstart "./../../source/runcstart"
	specs "./../../source/specs"
	"fmt"
	"log"
	"strings"
)

func setupConfig(process Process, configFile string, rootPath string) error {
	var linuxspec *specs.LinuxSpec
	linuxspec, err := configconvert.ConfigToLinuxSpec(configFile)
	if err != nil {
		log.Fatalf("[Specstestroot] Process readconfig error, err = %v", err)
	}

	linuxspec.Spec.Process.User.UID = process.User.UID
	linuxspec.Spec.Process.User.GID = process.User.GID
	linuxspec.Spec.Process.User.AdditionalGids = process.User.AdditionalGids
	linuxspec.Spec.Process.Env = process.Env
	//linuxspec.Spec.Process.User.GID = process.Env[1]
	linuxspec.Spec.Process.Cwd = process.Cwd
	linuxspec.Spec.Process.Args = process.Args

	linuxspec.Spec.Root.Path = rootPath
	err = configconvert.LinuxSpecToConfig(configFile, linuxspec)
	return err
}
func testProcess(process Process, rootPath string) (bool, error) {
	configFile := "./../../source/config.json"
	err := setupConfig(process, configFile, rootPath)
	if err != nil {
		log.Fatalf("[Specstest] process = %v setupEnv failed, err = %v...", process, err)
	} else {
		fmt.Printf("[Specstest] process = %v setupEnv sucess ... \n", process)
	}

	var retBool bool
	output, err := runcstart.StartRunc(configFile)
	if err != nil {
		retBool = false
	} else {
		if strings.Contains(output, "runc ran") {
			retBool = true
		} else {
			retBool = false
		}
	}

	return retBool, err
}

func setResult(testResult *TestResult, key []string, res string) {
	testResult.Process.Terminal[key[0]] = res
	testResult.Process.User.UID[key[1]] = res
	testResult.Process.User.GID[key[2]] = res
	testResult.Process.User.AdditionalGids[key[3]] = res
	testResult.Process.Env[key[4]] = res
	testResult.Process.Env[key[5]] = res
	testResult.Process.Cwd[key[6]] = res
	testResult.Process.Args[key[7]] = res
}
