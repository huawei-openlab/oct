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
	hostsetup "./../../source/hostsetup"
	specs "./../../source/specs"

	"fmt"
	"log"
	"os"
	"os/exec"
)

func setupEnv(testValue string) error {
	programmeString := ""
	outputFile := ""
	//testFile := "test_version_correct"
	err := hostsetup.SetupEnv(programmeString, outputFile)
	if err != nil {
		return err
	}

	var filePath string
	filePath = "./../../source/config.json"
	var linuxspec *specs.LinuxSpec
	linuxspec, err = configconvert.ConfigToLinuxSpec(filePath)
	if err != nil {
		return err
	}

	linuxspec.Spec.Root.Path = "./../../source/rootfs_rootconfig"
	linuxspec.Spec.Version = testValue
	// Args take runc to run "ls" and then exit, we should not need "ls" here, but Args should not be nil
	linuxspec.Spec.Process.Args[0] = "ls"
	err = configconvert.LinuxSpecToConfig(filePath, linuxspec)
	return err
}

func startRunc(configFile string) error {
	cmd := exec.Command("runc", configFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

func testVersion(testValue string) error {
	err := setupEnv(testValue)
	if err != nil {
		log.Fatalf("[Specstest] version = %s setupEnv failed, err = %v...", testValue, err)
	} else {
		fmt.Println("[Specstest] version = %s setupEnv sucess ...", testValue)
	}
	configFile := "./../../source/config.json"
	err = startRunc(configFile)
	return err
}
