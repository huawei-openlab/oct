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
	"flag"
	"fmt"
	// _ "github.com/huawei-openlab/oct/tools/specstest/cases/linuxcapabilities"
	_ "github.com/huawei-openlab/oct/tools/specstest/cases/linuxnamespace"
	// _ "github.com/huawei-openlab/oct/tools/specstest/cases/linuxresources"
	// _ "github.com/huawei-openlab/oct/tools/specstest/cases/linuxrlimits"
	// _ "github.com/huawei-openlab/oct/tools/specstest/cases/linuxselinuxlabel"
	// _ "github.com/huawei-openlab/oct/tools/specstest/cases/linuxsysctl"
	_ "github.com/huawei-openlab/oct/tools/specstest/cases/specmount"
	_ "github.com/huawei-openlab/oct/tools/specstest/cases/specplatform"
	_ "github.com/huawei-openlab/oct/tools/specstest/cases/specroot"
	_ "github.com/huawei-openlab/oct/tools/specstest/cases/specversion"
	"github.com/huawei-openlab/oct/tools/specstest/hostenv"
	"github.com/huawei-openlab/oct/tools/specstest/manager"
	"github.com/huawei-openlab/oct/tools/specstest/utils"
	"log"
)

var specsRev = flag.String("specs", "", "Specify specs Revision from opencontainers/specs as the benchmark, in the form of commit id")
var runcRev = flag.String("runc", "", "Specify runc Revision from opencontainers/specs to be tested, in the form of commit id")

/*var output = flag.String("o", "./report/specs.json", "Specify filePath to install the test result, file with result should be in json format")*/

func main() {

	flag.Parse()

	if *specsRev == "" {
		fmt.Println("It is going to use newest specs revision as a test benchmark")
	}

	if *runcRev == "" {
		fmt.Println("It is going to test specs on newest revision")
	}

	if *specsRev != "" && *runcRev != "" {
		hostenv.UpateSpecsRev(*specsRev)
		hostenv.UpateRuncRev(*runcRev)
		fmt.Printf("It is going to do specs test on runc revesion commit %v specs revesion commit %v \n", *specsRev, *runcRev)
	}

	err := hostenv.SetupEnv("", "")
	if err != nil {
		log.Fatalf(" Pull image error, %v", err)
	}

	/*linuxnamespace.TestSuiteNP.Run()
	result := linuxnamespace.TestSuiteNP.GetResult()

	err = utils.StringOutput("namespace_out.json", result)
	if err != nil {
		log.Fatalf("Write namespace out file error,%v\n", err)
	}
	*/
	// // spec.version test
	/*specversion.TestSuiteVersion.Run()
	result = specversion.TestSuiteVersion.GetResult()

	err = utils.StringOutput("Version_out.json", result)
	if err != nil {
		log.Fatalf("Write version out file error,%v\n", err)
	}*/

	// // spec.mount test
	/*specmount.TestSuiteMount.Run()
	result = specmount.TestSuiteMount.GetResult()
	err = utils.StringOutput("Mount_out.json", result)
	if err != nil {
		log.Fatalf("Write mount out file error,%v\n", err)
	}*/
	// manager * TestManager = new(TestManager)

	for _, ts := range manager.Manager.TestSuite {
		ts.Run()
		result := ts.GetResult()
		outputJson := ts.Name + ".json"
		err := utils.StringOutput(outputJson, result)
		if err != nil {
			log.Fatalf("Write %v out file error,%v\n", ts.Name, err)
		}

	}
	/*
		specroot.TestSuiteRoot.Run()
		result := specroot.TestSuiteRoot.GetResult()
		err = utils.StringOutput("Root_out.json", result)
		if err != nil {
			log.Fatalf("Write Root out file error,%v\n", err)
		}*/

	/*specplatform.TestSuitePlatform.Run()
	result = specplatform.TestSuitePlatform.GetResult()
	err = utils.StringOutput("Platform_out.json", result)
	if err != nil {
		log.Fatalf("Write Platform out file error,%v\n", err)
	}
	*/
	//linux resources test
	/*linuxresources.TestSuiteLinuxResources.Run()
	result := linuxresources.TestSuiteLinuxResources.GetResult()
	err = ioutil.WriteFile("LinuxResources_out.json", []byte(result), 0777)
	if err != nil {
		log.Fatalf("Write LinuxResources out file error,%v\n", err)
	}*/

}
