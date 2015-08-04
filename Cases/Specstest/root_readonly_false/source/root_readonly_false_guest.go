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
	"log"
	"os"
	"os/exec"
	"strings"
)

func testRootReadonlyFalse() {

	//exec shell in host machine to get docker container id
	rootString := " / "
	cmd := exec.Command("/bin/sh", "-c", "mount |grep "+rootString)
	outBytes, err := cmd.Output()
	if err != nil {
		log.Fatalf("Specs test testRootReadonlyFalse grep mount string err, %v", err)
	}

	outString := string(outBytes)
	var resultString string
	if strings.Contains(outString, "(rw,") {
		resultString = "[YES]        Linuxspec.Spec.Root.Readonly == false   passed"
	} else {
		resultString = "[NO]        Linuxspec.Spec.Root.Readonly == false   failed"
	}

	foutfile := "/testtool/readonly_false_out.txt"
	fout, err := os.Create(foutfile)
	defer fout.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		fout.WriteString(resultString)
	}
}

func main() {
	testRootReadonlyFalse()
}
