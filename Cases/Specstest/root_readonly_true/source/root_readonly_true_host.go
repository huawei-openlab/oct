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
	"fmt"
	specs "github.com/opencontainers/specs"
	"log"
	"os/exec"
)

func testRootReadonlyTrue() {
	var hostpasswd string
	hostpasswd = "565513"

	cmd := exec.Command("/bin/sh", "-c", "echo "+hostpasswd+" | sudo mkdir -p /tmp/testtool")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: pull image error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "docker pull ubuntu:14.04")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: pull image error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "docker export $(docker create ubuntu) > ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: export image error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "mkdir -p rootfs_rootconfig")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: create rootfs dir error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "echo "+hostpasswd+" | sudo tar -C rootfs_rootconfig -xf ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: create rootfs content error, %v", err)
	}

	var filePath string
	filePath = "config.json"

	var linuxspec *specs.LinuxSpec
	linuxspec, err = configconvert.ConfigToLinuxSpec(filePath)
	if err != nil {
		log.Fatalf("Specstestroot readonly test: readconfig error, %v", err)
	}

	linuxspec.Spec.Root.Path = "rootfs_rootconfig"
	linuxspec.Spec.Root.Readonly = true
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
