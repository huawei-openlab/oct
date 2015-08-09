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
	//"bytes"
	//"fmt"
	//specs "github.com/opencontainers/specs"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
)

type TestResult struct {
	Readonly map[string]string `json:"Linuxspec.Spec.Root.Readonly"`
}

func testRootReadonlyTrue() {

	//exec shell in host machine to get docker container id
	rootString := " / "
	cmd := exec.Command("/bin/sh", "-c", "mount |grep "+rootString)
	outBytes, err := cmd.Output()
	if err != nil {
		log.Fatalf("Specs test testRootReadonlyTrue grep mount string err, %v", err)
	}

	testResult := new(TestResult)
	testResult.Readonly = make(map[string]string)
	outString := string(outBytes)
	if strings.Contains(outString, "(ro,") {
		testResult.Readonly["true"] = "passed"
	} else {
		testResult.Readonly["true"] = "failed"
	}

	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("Convert to json err, error:  %v\n", err)
		return
	}

	foutfile := "/testtool/readonly_true_out.json"
	fout, err := os.Create(foutfile)
	defer fout.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		fout.WriteString(string(jsonString))
	}
}

func main() {
	testRootReadonlyTrue()
}
