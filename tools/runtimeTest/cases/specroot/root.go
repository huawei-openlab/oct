// +build predraft

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
	"github.com/opencontainers/specs"
	"runtime"
	"strings"
)

/**
*Need mount proc and set mnt namespace when get namespace from container
*and the specs.Process.Terminal must be false when call runc in programe.
 */
var linuxSpec specs.LinuxSpec = specs.LinuxSpec{
	Spec: specs.Spec{
		Version: "pre-draft",
		Platform: specs.Platform{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		Root: specs.Root{
			Path:     "rootfs",
			Readonly: true,
		},
		Process: specs.Process{
			Terminal: false,
			User: specs.User{
				UID:            0,
				GID:            0,
				AdditionalGids: nil,
			},
			Args: []string{""},
		},
		Mounts: []specs.Mount{
			{
				Type:        "proc",
				Source:      "proc",
				Destination: "/proc",
				Options:     "",
			},
		},
	},
	Linux: specs.Linux{
		Resources: specs.Resources{
			Memory: specs.Memory{
				Swappiness: -1,
			},
		},
		Namespaces: []specs.Namespace{
			{
				Type: "mount",
				Path: "",
			},
		},
	},
}

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

	linuxSpec.Root.Readonly = readonlyValue
	linuxSpec.Root.Path = path
	return linuxSpec
}

func testRoot(linuxSpec *specs.LinuxSpec, readonlyValue bool, pathValue string) (string, error) {
	configFile := "./config.json"
	linuxSpec.Spec.Process.Args[0] = "/bin/mount"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)

	out, err := adaptor.StartRunc(configFile)
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
			return manager.FAILED, nil
		}
	} else {
		return manager.UNKNOWNERR, nil
	}
}
