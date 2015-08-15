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

	"fmt"
	"log"
	"os"
	"os/exec"
)

type TestResult struct {
	Capabilities map[string]string `json:"Linuxspec.Linux.Capabilities"`
}

func linuxCapabilitiesTestSETFCAP() {
	// init the output json file
	testResult := new(TestResult)
	testResult.Capabilities = make(map[string]string)

	//authorized  /testtool/linux_capabilities_SETFCAP_guest the capbility to  set the capbility of another file
	cmd := exec.Command("/bin/sh", "-c", "setcap CAP_SETFCAP=eip /testtool/linux_capabilities_SETFCAP_guest")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux Capabilities SETFCAP set the SETFCAP Capability error, %v", err)
		fmt.Println("[Specstest] linux Capabilities SETFCAP set the SETFCAP Capability error, %v", err)
		testResult.Capabilities["SETFCAP"] = "failed"
		fmt.Println("testResult.Capabilities[SETFCAP] = failed")
	} else {
		testResult.Capabilities["SETFCAP"] = "pass"
		fmt.Println("testResult.Capabilities[SETFCAP] = pass")
	}

	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("Convert to json err, error:  %v\n", err)
		fmt.Println("Convert to json err, error:  %v\n", err)
		return
	}

	foutfile := "/testtool/linux_capabilities_SETFCAP.json"
	fout, err := os.Create(foutfile)
	defer fout.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		fout.WriteString(string(jsonString))
	}
}

func main() {
	linuxCapabilitiesTestSETFCAP()
}
