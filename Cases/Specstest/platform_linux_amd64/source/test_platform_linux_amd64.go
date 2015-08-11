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
	Platform PlatformStr `json:"Linuxspec.Spec.Platform"`
}

type PlatformStr struct {
	OS   map[string]string `json:"OS"`
	Arch map[string]string `json: "Arch"`
}

func testPlatformLinuxAmd64() {

	testResult := new(TestResult)
	testResult.Platform.OS = make(map[string]string)
	testResult.Platform.Arch = make(map[string]string)

	cmd := exec.Command("runc")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		testResult.Platform.OS["linux"] = "failed"
		testResult.Platform.Arch["amd64"] = "failed"
	} else {
		testResult.Platform.OS["linux"] = "passed"
		testResult.Platform.Arch["amd64"] = "passed"
	}

	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("Convert to json err, error:  %v\n", err)
		return
	}

	foutfile := "/tmp/testtool/platform_linux_amd64_out.json"
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
