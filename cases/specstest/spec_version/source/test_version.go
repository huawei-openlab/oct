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
	specs "./../../source/specs"
	"encoding/json"
	"fmt"
	"log"
	"os"
	//"os/exec"
)

type TestResult struct {
	Version map[string]string `json:"Linuxspec.Spec.Version"`
}

func main() {
	testResult := new(TestResult)
	testResult.Version = make(map[string]string)

	passString := "passed"
	failedString := "failed"
	//Do version test when value is correct
	testValue := specs.Version
	err := testVersion(testValue)
	if err != nil {
		fmt.Printf("[Specstest] Version = %s testVersion err = %v ... \n", testValue, err)
		testResult.Version[testValue] = failedString
	} else {
		testResult.Version[testValue] = passString
	}

	//Do version test when value is err
	testValue = specs.Version + "_err"
	err = testVersion(testValue)
	if err != nil {
		fmt.Printf("[Specstest] Version = %s testVersion err = %v ... \n", testValue, err)
		testResult.Version[testValue] = passString
	} else {
		testResult.Version[testValue] = failedString
	}

	// Write result to ouput json file.
	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("[Specstest] testResult = %v convert to json err err = %v\n", testResult, err)
		return
	}

	foutfile := "/tmp/testtool/spec_version.json"
	fout, err := os.Create(foutfile)
	defer fout.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		fout.WriteString(string(jsonString))
	}
}
