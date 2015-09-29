// Copyright 2015 The oct Authors
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
	"./sc"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/appc/spec/schema"
	"github.com/opencontainers/specs"
	"os"
)

func convertRocketFile(path string) {
	var image schema.ImageManifest
	var ls specs.LinuxSpec
	var msgs []string

	content, err := ReadFile(path)
	if err != nil {
		return
	}
	json.Unmarshal([]byte(content), &image)

	ls, msgs = specsConvert.LinuxSpecFrom(image, msgs)

	val, _ := json.MarshalIndent(ls, "", "\t")
	fmt.Println(string(val))

}

func ReadFile(file_url string) (content string, err error) {
	_, err = os.Stat(file_url)
	if err != nil {
		fmt.Println("cannot find the file ", file_url)
		return content, err
	}
	file, err := os.Open(file_url)
	defer file.Close()
	if err != nil {
		fmt.Println("cannot open the file ", file_url)
		return content, err
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	content = buf.String()

	return content, nil
}
