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

func setupConfig(testValue bool, configFile string, rootPath string) error {
	var linuxspec *specs.LinuxSpec
	linuxspec, err := configconvert.ConfigToLinuxSpec(configFile)
	if err != nil {
		log.Fatalf("[Specstestroot] Root test readonly = %v readconfig error, err = %v", testValue, err)
	}

	linuxspec.Spec.Root.Path = rootPath
	linuxspec.Spec.Root.Readonly = testValue
	linuxspec.Spec.Process.Args[0] = "mount"
	err = configconvert.LinuxSpecToConfig(configFile, linuxspec)
	return err
}

func testRoot(testValue bool, rootPath string) (bool, error) {
	configFile := "./../../source/config.json"
	err := setupConfig(testValue, configFile, rootPath)
	var retString bool
	if err != nil {
		log.Fatalf("[Specstest] root.readonly = %v setupEnv failed, err = %v...", testValue, err)
	} else {
		fmt.Printf("[Specstest] root.readonly = %v setupEnv sucess ... \n", testValue)
	}

	output, err := runcstart.StartRunc(configFile)
	if err != nil {
		log.Fatalf("[Specstest] root.readonly = %v startRunc failed, err = %v...", testValue, err)
	}
	if testValue == true && strings.Contains(output, "(ro,") {
		retString = true
	} else if testValue == false && strings.Contains(output, "(rw,") {
		retString = true
	} else {
		retString = false
	}

	return retString, err
}
