// +build v0.1.1

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

package specversion

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
)

const (
	testValueCorrect = "0.1.0"
	testVauleError   = "0.1.0-err"
)

// Init TestSuite for spec.Version
var TestSuiteVersion manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Version"}

// Add TestCase to Testsuite, add TestSuite to TestManager
func init() {
	TestSuiteVersion.AddTestCase("TestVersionCorrect", TestVersionCorrect)
	TestSuiteVersion.AddTestCase("TestVersionError", TestVersionError)
	manager.Manager.AddTestSuite(TestSuiteVersion)
}

// Set input value of spec.Version to specs.LinuxSpec obj
func setVersion(testValue string) specs.LinuxSpec {
	// Get smallest set of specs.LinuxSpec
	ls := specsinit.SetLinuxspecMinimum()
	// Set value
	ls.Version = testValue
	return ls
}

// Convert specs obj to json file, and start runc to have a test, return testresult
func testVersion(linuxSpec *specs.LinuxSpec, linuxRuntime *specs.LinuxRuntimeSpec, value bool) (string, error) {

	// Convert LinuxSpec to config.json
	configFile := "./config.json"
	linuxSpec.Spec.Process.Args[0] = "/bin/ls"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	if err != nil {
		return manager.UNKNOWNERR, err
	}

	// Convert LinuxRuntimeSpec to runtime.json
	runtimeFile := "./runtime.json"
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntime)
	if err != nil {
		return manager.UNKNOWNERR, err
	}

	// Start runc to have a test
	out, err := adaptor.StartRunc(configFile, runtimeFile)
	if err != nil {
		if value {
			return manager.FAILED, errors.New(string(out) + err.Error())
		} else {
			return manager.PASSED, nil
		}

	} else {
		if value {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New("Give wrong value to Version but runc works!")
		}
	}
}
