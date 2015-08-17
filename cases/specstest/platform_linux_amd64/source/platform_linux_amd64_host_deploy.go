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
	"os/exec"
)

func testPlatform() {
	programmeString := "demo"
	outputFile := "platform_linux_amd64_out"
	err := hostsetup.SetupEnv(programmeString, outputFile)
	if err != nil {
		log.Fatalf("[Specstest] platform linux amd64 test: hostsetup.SetupEnv error, %v", err)
	}

	fmt.Println("Build test programme...................")
	cmd := exec.Command("/bin/sh", "-c", "go build test_platform_linux_amd64.go")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] platform linux amd64 test: build test programme error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "mv test_platform_linux_amd64 ./../../source")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] platform linux amd64 test:: mv test programme error, %v", err)
	}

	fmt.Println("Host enviroment setting up for runc is already!")
	var filePath string
	filePath = "./../../source/config.json"

	oriSpecVersion := specs.Version
	var linuxspec *specs.LinuxSpec
	linuxspec, err = configconvert.ConfigToLinuxSpec(filePath)
	if err != nil {
		log.Fatalf("[Specstest] platform linux amd64 test: readconfig error, %v", err)
	}

	linuxspec.Spec.Root.Path = "./rootfs_rootconfig"
	linuxspec.Spec.Version = oriSpecVersion
	linuxspec.Spec.Platform.OS = "linux"
	linuxspec.Spec.Platform.Arch = "amd64"
	linuxspec.Spec.Process.Args[0] = "./" + programmeString
	err = configconvert.LinuxSpecToConfig(filePath, linuxspec)
	if err != nil {
		log.Fatalf("[Specstest] platform linux amd64 test: writeconfig error, %v", err)
	}
	fmt.Println("Host enviroment for runc is already!")

}

func main() {
	testPlatform()
}
