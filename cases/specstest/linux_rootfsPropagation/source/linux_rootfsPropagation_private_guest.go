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
	"log"
	"os/exec"
)

type TestResult struct {
	RootfsPropagation map[string]string `json:"Linuxspec.Linux.RootfsPropagation"`
}

func rootfsPropagationTestPrivate() {

	//touch a new file and test whether it can be seen from the host machine
	cmd := exec.Command("/bin/sh", "-c", "echo \"YouCanSeeThisContent\" > /rootfsPropagationTestPrivate.txt")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux rootfsPropagation private test touch the file error, %v", err)
	}
}

func main() {
	rootfsPropagationTestPrivate()
}
