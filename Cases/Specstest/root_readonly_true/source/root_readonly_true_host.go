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

	cmd := exec.Command("/bin/sh", "-c", "go build root_readonly_true_guest.go")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: build guest programme error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "mkdir -p /tmp/testtool")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: mkdir testtool dir error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "mv root_readonly_true_guest /tmp/testtool/")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: mv guest programme error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "touch  /tmp/testtool/readonly_true_out.txt")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: touch readonly_true_out.txt err, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "chown root:root  /tmp/testtool/readonly_true_out.txt")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: change to root power err, %v", err)
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

	cmd = exec.Command("/bin/sh", "-c", "mkdir -p ./../../source/rootfs_rootconfig")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: create rootfs dir error, %v", err)
	}

	cmd = exec.Command("/bin/sh", "-c", "tar -C ./../../source/rootfs_rootconfig -xf ubuntu.tar")
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

	linuxspec.Spec.Root.Path = "./../../source/rootfs_rootconfig"
	linuxspec.Spec.Root.Readonly = true
	err = configconvert.LinuxSpecToConfig(filePath, linuxspec)
	//err = wirteConfig(filePath, linuxspec)
	if err != nil {
		log.Fatalf("Specstest root readonly test: writeconfig error, %v", err)
	}
	fmt.Println("Host enviroment for runc is already!")

	/*
		cmd = exec.Command("/bin/bash", "-c", "runc")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: start runc and test error, %v", err)
		}*/
}

func main() {
	testRootReadonlyTrue()
}
