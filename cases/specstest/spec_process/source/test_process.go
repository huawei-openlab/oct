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
	//specs "./../../source/specs"
	"encoding/json"
	"log"
	//"os"
	//"os/exec"
	hostsetup "./../../source/hostsetup"
	"fmt"
)

type TestResult struct {
	Process ProcessStr `json:"Linuxspec.Spec.Process"`
}

type ProcessStr struct {
	// Terminal creates an interactive terminal for the container.
	Terminal map[string]string `json:"terminal"`
	// User specifies user information for the process.
	User UserStr `json:"user"`
	// Args specifies the binary and arguments for the application to execute.
	Args map[string]string `json:"args"`
	// Env populates the process environment for the process.
	Env map[string]string `json:"env"`
	// Cwd is the current working directory for the process and must be
	// relative to the container's root.
	Cwd map[string]string `json:"cwd"`
}

type UserStr struct {
	UID            map[string]string `json:"uid"`
	GID            map[string]string `json:"gid"`
	AdditionalGids map[string]string `json:"additionalGids"`
}

type Process struct {
	// Terminal creates an interactive terminal for the container.
	Terminal bool `json:"terminal"`
	// User specifies user information for the process.
	User User `json:"user"`
	// Args specifies the binary and arguments for the application to execute.
	Args []string `json:"args"`
	// Env populates the process environment for the process.
	Env []string `json:"env"`
	// Cwd is the current working directory for the process and must be
	// relative to the container's root.
	Cwd string `json:"cwd"`
}

// User specifies Linux specific user and group information for the container's
// main process
type User struct {
	// Uid is the user id
	UID int32 `json:"uid"`
	// Gid is the group id
	GID int32 `json:"gid"`
	// AdditionalGids are additional group ids set for the container's process
	AdditionalGids []int32 `json:"additionalGids"`
}

func main() {
	testResult := new(TestResult)
	testResult.Process.Terminal = make(map[string]string)
	testResult.Process.Args = make(map[string]string)
	testResult.Process.Env = make(map[string]string)
	testResult.Process.Cwd = make(map[string]string)

	testResult.Process.User.UID = make(map[string]string)
	testResult.Process.User.GID = make(map[string]string)
	testResult.Process.User.AdditionalGids = make(map[string]string)

	outputFile := "spec_process"
	guestProgrammeFile := "process_guest"
	err := hostsetup.SetupEnv(guestProgrammeFile, outputFile)
	if err != nil {
		log.Fatalf("[Specstest] process test: hostsetup.SetupEnv error, %v", err)
	}
	fmt.Println("Host enviroment setting up for runc is already!")

	user := User{
		UID:            1,
		GID:            1,
		AdditionalGids: nil}
	process := Process{
		Terminal: false,
		User:     user,
		Env:      []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", "TERM=xterm"},
		Cwd:      "/testtool",
		Args:     []string{"./process_guest"}}

	var key = []string{"false", "1", "1", "nil", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"TERM=xterm", "/testtool", "./process_guest"}

	rootPath := "./../../source/rootfs_rootconfig"
	retBool, err := testProcess(process, rootPath)
	pString := "passed"
	fString := "failed"
	npString := "notSupported"
	if err != nil {
		setResult(testResult, key, npString)
	} else {
		if retBool == true {
			setResult(testResult, key, pString)
		} else {
			setResult(testResult, key, fString)
		}
	}

	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("[Specstest] Process test ouput json failed, err = %v...", err)
	}
	hostsetup.HostOutput(outputFile, string(jsonString))
}
