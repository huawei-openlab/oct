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
	"github.com/huawei-openlab/oct/tools/specsValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/specsValidator/manager"
	"github.com/huawei-openlab/oct/tools/specsValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/specsValidator/utils/specsinit"
	"github.com/opencontainers/specs"
)

const (
	testValueCorrect = "0.1.0"
	testVauleError   = "0.1.0-err"
)

var TestSuiteVersion manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Version"}

func init() {
	TestSuiteVersion.AddTestCase("TestVersionCorrect", TestVersionCorrect)
	TestSuiteVersion.AddTestCase("TestVersionError", TestVersionError)
	manager.Manager.AddTestSuite(TestSuiteVersion)
}

func setVersion(testValue string) specs.LinuxSpec {
	ls := specsinit.SetLinuxspecMinimum()
	ls.Version = testValue
	return ls
}

func testVersion(linuxSpec *specs.LinuxSpec, linuxRuntime *specs.LinuxRuntimeSpec, value bool) (string, error) {
	configFile := "./config.json"
	linuxSpec.Spec.Process.Args[0] = "/bin/ls"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	if err != nil {
		return manager.UNKNOWNERR, err
	}

	runtimeFile := "./runtime.json"
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntime)
	if err != nil {
		return manager.UNKNOWNERR, err
	}

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
