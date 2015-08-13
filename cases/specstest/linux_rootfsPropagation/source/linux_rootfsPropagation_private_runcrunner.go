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
	"strings"
)

type TestResult struct {
	RootfsPropagation map[string]string `json:"Linuxspec.Linux.RootfsPropagation"`
}

func rootfsPropagationTestPrivate() {
	//excute the runc
	cmd := exec.Command("runc")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("Start: %v", err)
		fmt.Println("Start: %v", err)
	}

	// init the output json file
	testResult := new(TestResult)
	testResult.RootfsPropagation = make(map[string]string)

	//cat the file touched inside the container verify whether the mount propagation works
	cmd1 := exec.Command("cat", "./rootfs_rootconfig/rootfsPropagationTestPrivate.txt")
	cmdouput, err1 := cmd1.Output()
	var comparestring, cmdout string
	comparestring = "YouCanSeeThisContent"
	cmdout = strings.TrimSpace(string(cmdouput))
	if err1 != nil {
		log.Fatalf("[Specstest] linux rootfsPropagation private test read the testfileerror, %v", err1)
		testResult.RootfsPropagation["private"] = "pass"
	} else {
		if strings.EqualFold(cmdout, comparestring) {
			testResult.RootfsPropagation["private"] = "failed"
		}
	}

	//output the json file
	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("Convert to json err, error:  %v\n", err)
		return
	}

	foutfile := "/tmp/testtool/linux_rootfsPropagation_private.json"
	fout, err := os.Create(foutfile)
	defer fout.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		fout.WriteString(string(jsonString))
	}
}

func main() {
	rootfsPropagationTestPrivate()
}
