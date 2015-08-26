// Copyright 2015 Huawei Inc. All Rights Reserved.
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
	"io/ioutil"
	"log"

	"github.com/huawei-openlab/oct/cases/specstest/cases/specplatform"
	"github.com/huawei-openlab/oct/cases/specstest/cases/specroot"
	"github.com/huawei-openlab/oct/cases/specstest/cases/linuxnamespace"
	"github.com/huawei-openlab/oct/cases/specstest/cases/specmount"
	"github.com/huawei-openlab/oct/cases/specstest/cases/specversion"
	"github.com/huawei-openlab/oct/cases/specstest/hostenv"
)

func main() {

	err := hostenv.SetupEnv("", "")
	if err != nil {
		log.Fatalf(" Pull image error, %v", err)
	}
	linuxnamespace.TestSuiteNP.Run()
	result := linuxnamespace.TestSuiteNP.GetResult()

	err = ioutil.WriteFile("namespace_out.json", []byte(result), 0777)
	if err != nil {
		log.Fatalf("Write namespace out file error,%v\n", err)
	}

	// spec.version test
	specversion.TestSuiteVersion.Run()
	result = specversion.TestSuiteVersion.GetResult()

	err = ioutil.WriteFile("Version_out.json", []byte(result), 0777)
	if err != nil {
		log.Fatalf("Write version out file error,%v\n", err)
	}

	// spec.mount test
	specmount.TestSuiteMount.Run()
	result = specmount.TestSuiteMount.GetResult()
	err = ioutil.WriteFile("Mount_out.json", []byte(result), 0777)
	if err != nil {
		log.Fatalf("Write mount out file error,%v\n", err)
	}

	specroot.TestSuiteRoot.Run()
	result = specroot.TestSuiteRoot.GetResult()

	err = ioutil.WriteFile("Root_out.json", []byte(result), 0777)
	if err != nil {
		log.Fatalf("Write Root out file error,%v\n", err)
	}

	specplatform.TestSuitePlatform.Run()
	result = specplatform.TestSuitePlatform.GetResult()
	err = ioutil.WriteFile("Platform_out.json", []byte(result), 0777)
	if err != nil {
		log.Fatalf("Write Platform out file error,%v\n", err)
	}

}
