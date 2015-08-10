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
	"encoding/json"
	//"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type TestResult struct {
	Readonly map[string]string `json:"Linuxspec.Spec.Root.Readonly"`
}

func testRootReadonlyFalse() {

	//exec shell in host machine to get docker container id
	rootString := " / "
	cmd := exec.Command("/bin/sh", "-c", "mount |grep "+rootString)
	outBytes, err := cmd.Output()
	if err != nil {
		log.Fatalf("Specs test testRootReadonlyFalse grep mount string err, %v", err)
	}
	testResult := new(TestResult)
	testResult.Readonly = make(map[string]string)
	outString := string(outBytes)
	//var resultString string
	if strings.Contains(outString, "(rw,") {
		testResult.Readonly["false"] = "passed"
		//resultString = "[YES]        Linuxspec.Spec.Root.Readonly == false   passed"
	} else {
		testResult.Readonly["false"] = "failed"
		//resultString = "[NO]        Linuxspec.Spec.Root.Readonly == false   failed"
	}

	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("Convert to json err, error:  %v\n", err)
		return
	}

	foutfile := "/testtool/readonly_false_out.json"
	//fout, err := os.OpenFile(foutfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	fout, err := os.Create(foutfile)
	//defer fout.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		/*
			var b []byte
			n, err := fout.Read(b)
			if err != nil {
				log.Fatalf("Read file err : %v", err)

			} else {
				fmt.Println("len n : %v", n)
				fout.WriteAt(jsonString, int64(n))
			}
		*/
		fout.WriteString(string(jsonString))
	}
}

func main() {
	testRootReadonlyFalse()
}
