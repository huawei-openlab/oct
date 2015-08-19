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
	//configconvert "./../../source/configconvert"
	hostsetup "./../../source/hostsetup"
	"fmt"
	// specs "github.com/opencontainers/specs" //newest version
	//specs "./../../source/specs"
	"encoding/json"
	"log"
	"os"
)

type TestResult struct {
	Root RootStr `json:"Linuxspec.Spec.Root"`
}

type RootStr struct {
	// Path is the absolute path to the container's root filesystem.
	Path map[string]string `json:"path"`
	// Readonly makes the root filesystem for the container readonly before the process is executed.
	Readonly map[string]string `json:"readonly"`
}

func main() {
	guestProgrammeFile := ""
	outputFile := "spec_root"
	err := hostsetup.SetupEnv(guestProgrammeFile, outputFile)
	if err != nil {
		log.Fatalf("[Specstest] root test: hostsetup.SetupEnv error, %v", err)
	}
	fmt.Println("Host enviroment setting up for runc is already!")

	testResult := new(TestResult)
	testResult.Root.Path = make(map[string]string)
	testResult.Root.Readonly = make(map[string]string)

	//Start readonly == true test
	testValue := true
	rootPath := "./../../source/rootfs_rootconfig"
	retBool, err := testRoot(testValue, rootPath)
	if err != nil {
		log.Fatalf("[Specstest] root.readonly = %v testRoot failed, err = %v...", testValue, err)
	}
	if retBool == true {
		testResult.Root.Path[rootPath] = "passed"
		testResult.Root.Readonly["true"] = "passed"
	} else {
		testResult.Root.Path[rootPath] = "failed"
		testResult.Root.Readonly["false"] = "failed"
	}

	// Start readonly == false test
	testValue = false
	retBool, err = testRoot(testValue, rootPath)
	if err != nil {
		log.Fatalf("[Specstest] root.readonly = %v testRoot failed, err = %v...", testValue, err)
	}
	if retBool == true {
		testResult.Root.Path[rootPath] = "passed"
		testResult.Root.Readonly["false"] = "passed"
	} else {
		testResult.Root.Path[rootPath] = "failed"
		testResult.Root.Readonly["false"] = "failed"
	}

	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("[Specstest] root test ouput json failed, err = %v...", err)
	}

	foutfile := "/tmp/testtool/spec_root.json"
	fout, err := os.Create(foutfile)
	defer fout.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		fout.WriteString(string(jsonString))
	}
}
