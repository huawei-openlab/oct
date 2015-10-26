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

package specroot

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"strings"
)

const (
	testPathCorrect = "rootfs"
	testPathError   = "rootfs_error"
)

var TestSuiteRoot manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Root"}

func init() {
	TestSuiteRoot.AddTestCase("TestReadOnlyTrue", TestReadOnlyTrue)
	TestSuiteRoot.AddTestCase("TestReadOnlyFalse", TestReadOnlyFalse)
	TestSuiteRoot.AddTestCase("TestPathError", TestPathError)
	manager.Manager.AddTestSuite(TestSuiteRoot)
}

func setRoot(readonlyValue bool, path string) specs.LinuxSpec {
	ls := specsinit.SetLinuxspecMinimum()
	ls.Root.Readonly = readonlyValue
	ls.Root.Path = path
	return ls
}

func testRoot(linuxSpec *specs.LinuxSpec, linuxRuntime *specs.LinuxRuntimeSpec, readonlyValue bool, pathValue string) (string, error) {
	configFile := "./config.json"
	linuxSpec.Spec.Process.Args[0] = "/bin/mount"
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
		if pathValue == testPathError {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New(string(out) + err.Error())
		}
	}
	if pathValue == testPathCorrect {
		if readonlyValue == true && strings.Contains(out, "(ro,") {
			return manager.PASSED, nil
		} else if readonlyValue == false && strings.Contains(out, "(rw,") {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New("Start runc successful, but get the wrong right of root.")
		}
	} else {
		return manager.UNKNOWNERR, nil
	}
}
