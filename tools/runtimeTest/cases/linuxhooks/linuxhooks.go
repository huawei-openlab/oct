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

/*
// Hook specifies a command that is run at a particular event in the lifecycle of a container
type Hook struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
	Env  []string `json:"env"`
}

// Hooks for container setup and teardown
type Hooks struct {
	// Prestart is a list of hooks to be run before the container process is executed.
	// On Linux, they are run after the container namespaces are created.
	Prestart []Hook `json:"prestart"`
	// Poststop is a list of hooks to be run after the container process exits.
	Poststop []Hook `json:"poststop"`
}
*/

package linuxhooks

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
			Readonly: false,
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

var TestSuiteLinuxHooks manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Hooks"}

func init() {
	TestSuiteLinuxHooks.AddTestCase("TestSuiteLinuxHooksPrestart", TestSuiteLinuxHooksPrestart)
	// TestSuiteLinuxHooks.AddTestCase("TestSuiteLinuxHooksPoststop", TestSuiteLinuxHooksPoststop)
	manager.Manager.AddTestSuite(TestSuiteLinuxHooks)
}

func setHooks(thooks []specs.Hook, isPre bool) specs.LinuxSpec {
	if isPre {
		linuxSpec.Hooks.Prestart = thooks
	} else {
		linuxSpec.Hooks.Poststop = thooks
	}

	return linuxSpec
}

func testHooks(linuxSpec *specs.LinuxSpec, compare string, failinfo string) (string, error) {
	configFile := "./config.json"
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "ls"}
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	out, err := adaptor.StartRunc(configFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New(string(out) + err.Error())
	} else {
		if strings.Contains(strings.TrimSpace(string(out)), compare) {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New("test failed because" + failinfo)
		}
	}
}
