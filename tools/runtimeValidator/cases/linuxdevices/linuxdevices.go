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

package linuxdevices

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/opencontainers/specs"
	"runtime"
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
				Type:        "tmpfs",
				Source:      "tmpfs",
				Destination: "/dev",
				Options:     "nosuid,strictatime,mode=755,size=65536k",
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

var TestSuiteLinuxDevices manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Devices"}

// TODO
func init() {
	TestSuiteLinuxDevices.AddTestCase("TestSuiteLinuxDevicesFull", TestSuiteLinuxDevicesFull)
	manager.Manager.AddTestSuite(TestSuiteLinuxDevices)
}

func setDevices(testdevices specs.Device) specs.LinuxSpec {
	var initdevice specs.Device = specs.Device{
		Type:        99,
		Path:        "/dev/null",
		Major:       1,
		Minor:       3,
		Permissions: "rwm",
		FileMode:    438,
		UID:         0,
		GID:         0,
	}
	linuxSpec.Linux.Devices = []specs.Device{initdevice}
	linuxSpec.Linux.Devices = append(linuxSpec.Linux.Devices, testdevices)
	return linuxSpec
}
