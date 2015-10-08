// +build v0.1.1

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
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
)

/**
*Need mount proc and set mnt namespace when get namespace from container
*and the specs.Process.Terminal must be false when call runc in programe.
 */

var TestLinuxDevices manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Devices"}

// TODO
func init() {
	TestLinuxDevices.AddTestCase("TestLinuxDevicesFull", TestLinuxDevicesFull)
	manager.Manager.AddTestSuite(TestLinuxDevices)
}

func setDevices(testdevices specs.Device) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxRuntimeSpec := specsinit.SetLinuxruntimeMinimum()
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
	linuxRuntimeSpec.Linux.Devices = []specs.Device{initdevice}
	linuxRuntimeSpec.Linux.Devices = append(linuxRuntimeSpec.Linux.Devices, testdevices)
	return linuxSpec, linuxRuntimeSpec
}
