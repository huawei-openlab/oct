// +build predraft

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
	// _ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/linuxcapabilities"
	_ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/linuxnamespace"
	// _ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/linuxresources"
	// _ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/linuxrlimits"
	// _ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/linuxselinuxlabel"
	// _ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/linuxsysctl"
	_ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/specmount"
	_ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/specplatform"
	_ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/specprocess"
	_ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/specroot"
	_ "github.com/huawei-openlab/oct/tools/runtimeValidator/cases/specversion"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/hostenv"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"log"
)

var specsRev = flag.String("specs", "", "Specify specs Revision from opencontainers/specs as the benchmark, in the form of commit id")
var runcRev = flag.String("runc", "", "Specify runc Revision from opencontainers/specs to be tested, in the form of commit id")
var output = flag.String("o", "./report/", "Specify filePath to install the test result linuxspec.json")

func main() {

	flag.Parse()
	fmt.Println("Tested Revision:")
	if *specsRev == "" {
		defaultRev := "45ae53d4dba8e550942f7384914206103f6d2216"
		fmt.Printf("	Default specs revision: pre-draft(commit: %v)\n", defaultRev)
		hostenv.UpateSpecsRev(defaultRev)
	}

	if *runcRev == "" {
		defaultRev := "v0.0.4"
		fmt.Printf("	Default runc revision :  %v\n", defaultRev)
		hostenv.UpateRuncRev(defaultRev)
	}

	if *specsRev != "" && *runcRev != "" {
		hostenv.UpateSpecsRev(*specsRev)
		hostenv.UpateRuncRev(*runcRev)
		fmt.Printf("	Specified specs revision: %v \n", *specsRev)
		fmt.Printf("	Specified runc revision: %v \n", *runcRev)
	}

	if *output == "" {
		*output = "./report/"
		fmt.Println("Using default output: oct/tools/runtimeValidator/report/linuxspec.json")
	} else {
		fmt.Printf("Specified output: %v\n", *output)
	}

	err := hostenv.CreateBoundle()
	if err != nil {
		log.Fatalf("Create boundle error, %v", err)
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
	result := manager.Manager.GetTotalResult()
	err = utils.SpecifyOutput(*output, result)
	if err != nil {
		log.Fatalf("Write %v out file error,%v\n", *output, err)
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
