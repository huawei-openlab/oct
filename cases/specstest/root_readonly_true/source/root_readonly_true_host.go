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
	// specs "github.com/opencontainers/specs" //newest version
	specs "./../../source/specs"
	"log"
)

func testRootReadonlyTrue() {

	var guestProgrammeFile string
	guestProgrammeFile = "root_readonly_true_guest"
	var outputFile string
	outputFile = "readonly_true_out"
	err := hostsetup.SetupEnv(guestProgrammeFile, outputFile)
	if err != nil {
		log.Fatalf("Specstest root readonly test: hostsetup.SetupEnv error, %v", err)
	}
	fmt.Println("Host enviroment setting up for runc is already!")

	var filePath string
	filePath = "./../../source/config.json"

	var linuxspec *specs.LinuxSpec
	linuxspec, err = configconvert.ConfigToLinuxSpec(filePath)
	if err != nil {
		log.Fatalf("Specstestroot readonly test: readconfig error, %v", err)
	}

	linuxspec.Spec.Root.Path = "./rootfs_rootconfig"
	linuxspec.Spec.Root.Readonly = true
	linuxspec.Spec.Process.Args[0] = "./root_readonly_true_guest"
	linuxspec.Spec.Process.Terminal = true
	linuxspec.Spec.Process.User.UID = 0
	linuxspec.Spec.Process.User.GID = 0
	linuxspec.Spec.Process.User.AdditionalGids = nil
	err = configconvert.LinuxSpecToConfig(filePath, linuxspec)
	//err = wirteConfig(filePath, linuxspec)
	if err != nil {
		log.Fatalf("Specstest root readonly test: writeconfig error, %v", err)
	}
	fmt.Println("Host enviroment for runc is already!")

}

func main() {
	testRootReadonlyTrue()
}
