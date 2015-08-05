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
	"fmt"
	specs "github.com/opencontainers/specs"
	"log"
)

func testVersion() {
	programmeString := "demo"
	outputFile := ""
	err := hostsetup.SetupEnv(programmeString, outputFile)
	if err != nil {
		log.Fatalf("Specstest version test: hostsetup.SetupEnv error, %v", err)
	}
	fmt.Println("Host enviroment setting up for runc is already!")
	var filePath string
	filePath = "config.json"

	errSpecVersion := "0.1.0"
	var linuxspec *specs.LinuxSpec
	linuxspec, err = configconvert.ConfigToLinuxSpec(filePath)
	if err != nil {
		log.Fatalf("Specstest version test: readconfig error, %v", err)
	}

	linuxspec.Spec.Root.Path = "./../../source/rootfs_rootconfig"
	linuxspec.Spec.Version = errSpecVersion
	linuxspec.Spec.Process.Args[0] = "./" + programmeString
	err = configconvert.LinuxSpecToConfig(filePath, linuxspec)

	if err != nil {
		log.Fatalf("Specstest version test: writeconfig error, %v", err)
	}
	fmt.Println("Host enviroment for runc is already!")

}

func main() {
	testVersion()
}
