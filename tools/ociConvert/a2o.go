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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/appc/spec/schema"
	"github.com/huawei-openlab/oct/tools/ociConvert/lib"
	"github.com/opencontainers/specs"
	"os"
	"path"
)

func convertRocketFile(imagePath string, podPath string) {
	var image schema.ImageManifest
	var pod schema.PodManifest

	var ls specs.LinuxSpec
	var lrs specs.LinuxRuntimeSpec
	var msgs []string

	content, err := ReadFile(imagePath)
	if err != nil {
		fmt.Println("Cannot parse image file: ", imagePath)
		return
	}
	json.Unmarshal([]byte(content), &image)

	content, err = ReadFile(podPath)
	if err != nil {
		fmt.Println("Cannot parse pod file: ", podPath)
		return
	}
	json.Unmarshal([]byte(content), &pod)

	PreparePath("output", "")

	ls, msgs = specsConvert.LinuxSpecFrom(image, msgs)
	val, _ := json.MarshalIndent(ls, "", "\t")
	output := "output/config.json"
	fout, err := os.Create(output)
	if err != nil {
		fmt.Println(output, err)
	} else {
		fout.WriteString(string(val))
		fmt.Println("Generate ", output)
		fout.Close()
	}

	lrs, msgs = specsConvert.LinuxRuntimeSpecFrom(image, pod, msgs)
	val, _ = json.MarshalIndent(lrs, "", "\t")
	output = "output/runtime.json"
	fout, err = os.Create(output)
	if err != nil {
		fmt.Println(output, err)
	} else {
		fout.WriteString(string(val))
		fmt.Println("Generate ", output)
		fout.Close()
	}
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

func PreparePath(cachename string, filename string) (realurl string) {
	var dir string
	if filename == "" {
		dir = cachename
	} else {
		realurl = path.Join(cachename, filename)
		dir = path.Dir(realurl)
	}
	p, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(dir, 0777)
		}
	} else {
		if p.IsDir() {
			return realurl
		} else {
			os.Remove(dir)
			os.MkdirAll(dir, 0777)
		}
	}
	return realurl
}
