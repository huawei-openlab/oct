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
	"encoding/json"
	"log"
	"os"
	"os/exec"
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

func testPlatformLinuxAmd64() {

	testResult := new(TestResult)
	testResult.Process.Terminal = make(map[string]string)
	testResult.Process.Args = make(map[string]string)
	testResult.Process.Env = make(map[string]string)
	testResult.Process.Cwd = make(map[string]string)

	testResult.Process.User.UID = make(map[string]string)
	testResult.Process.User.GID = make(map[string]string)
	testResult.Process.User.AdditionalGids = make(map[string]string)

	cmd := exec.Command("runc")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {

		testResult.Process.Terminal["true"] = "failed"
		testResult.Process.User.UID["1"] = "failed"
		testResult.Process.User.GID["1"] = "failed"
		testResult.Process.User.AdditionalGids[""] = "failed"
		testResult.Process.Args["./process_guest"] = "failed"
		testResult.Process.Env["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"] = "failed"
		testResult.Process.Env["TERM=xterm"] = "failed"
		testResult.Process.Cwd["/testtool"] = "failed"
	} else {

		testResult.Process.Terminal["true"] = "passed"
		testResult.Process.User.UID["1"] = "passed"
		testResult.Process.User.UID["1"] = "passed"
		testResult.Process.User.AdditionalGids["nil"] = "passed"
		testResult.Process.Args["./process_guest"] = "passed"
		testResult.Process.Env["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"] = "passed"
		testResult.Process.Env["TERM=xterm"] = "passed"
		testResult.Process.Cwd["/testtool"] = "passed"
	}

	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("Convert to json err, error:  %v\n", err)
		return
	}

	foutfile := "/tmp/testtool/process_out.json"
	fout, err := os.Create(foutfile)
	defer fout.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		fout.WriteString(string(jsonString))
	}
}

func main() {
	testPlatformLinuxAmd64()
}
