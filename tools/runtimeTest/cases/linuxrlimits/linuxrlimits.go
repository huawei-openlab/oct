// +build predraft

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
//

package linuxrlimits

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
			{
				Type:        "cgroup",
				Source:      "cgroup",
				Destination: "/sys/fs/cgroup",
				Options:     "nosuid,noexec,nodev,relatime,ro",
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
			{
				Type: "pid",
				Path: "",
			},
			{
				Type: "network",
				Path: "",
			},
			{
				Type: "ipc",
				Path: "",
			},
			{
				Type: "uts",
				Path: "",
			},
		},
	},
}

var TestSuiteLinuxRlimits manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Rlimits"}

func init() {
	TestSuiteLinuxRlimits.AddTestCase("TestRlimitNOFILESoft", TestRlimitNOFILESoft)
	TestSuiteLinuxRlimits.AddTestCase("TestRlimitNOFILEHard", TestRlimitNOFILEHard)
	manager.Manager.AddTestSuite(TestSuiteLinuxRlimits)
}

func setRlimits(testrlimits specs.Rlimit) specs.LinuxSpec {
	linuxSpec.Linux.Rlimits = []specs.Rlimit{testrlimits}
	return linuxSpec
}

func testRlimits(linuxSpec *specs.LinuxSpec, rlimitItem string, value string, isSoftLimit bool) (string, error) {
	configFile := "./config.json"
	if isSoftLimit {
		linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "ulimit " + rlimitItem + " -S"}
	} else {
		linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "ulimit " + rlimitItem + " -H"}
	}

	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	out, err := adaptor.StartRunc(configFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New(string(out) + err.Error())
	} else {
		if strings.EqualFold(strings.TrimSpace(string(out)), value) {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, nil
		}
	}
}
