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

package configconvert

import (
	"encoding/json"
	"github.com/opencontainers/specs"
	"io/ioutil"
	"log"
	"os"
)

//read config.json to specs.LinuxSpec
func ConfigToLinuxSpec(filePath string) (*specs.LinuxSpec, error) {
	out, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var linuxspec specs.LinuxSpec
	err = json.Unmarshal(out, &linuxspec)
	if err != nil {
		return nil, err
	}

	return &linuxspec, nil
}

//write specs.LinuxSpec to config.json
func LinuxSpecToConfig(filePath string, linuxspec *specs.LinuxSpec) error {
	stream, err := json.Marshal(linuxspec)
	if err != nil {
		return err
	}
	objToJson(stream, filePath)
	return err
}

func objToJson(stream []byte, filePath string) {
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf(" open file err, %v", err)
	}
	defer fd.Close()
	_, err = fd.Write(stream)
	if err != nil {
		log.Fatalf(" write file err, %v", err)
	}
}
