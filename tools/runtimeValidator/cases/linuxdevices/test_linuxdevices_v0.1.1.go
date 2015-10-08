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
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/opencontainers/specs"
	"strings"
)

func TestLinuxDevicesFull() string {

	var device specs.Device = specs.Device{
		Type:        99,
		Path:        "/dev/full",
		Major:       1,
		Minor:       7,
		Permissions: "rwm",
		FileMode:    438,
		UID:         0,
		GID:         0,
	}
	linuxspec, linuxruntimespec := setDevices(device)
	utils.SetBind(&linuxruntimespec, &linuxspec)
	linuxspec.Spec.Process.Args[0] = "/containerend/linuxdevicesfull"

	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	err := configconvert.LinuxSpecToConfig(configFile, &linuxspec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, &linuxruntimespec)
	out, err := adaptor.StartRunc(configFile, runtimeFile)

	var result string
	var errout error
	if err != nil {
		result = manager.UNSPPORTED
		errout = errors.New("StartRunc error :" + out + ", " + err.Error())
	} else if strings.Contains(strings.TrimSpace(out), "echo: write error: No space left on device") {
		result = manager.PASSED
		errout = nil
	} else {
		result = manager.FAILED
		errout = errors.New("device /dev/full is NOT effective")
	}
	var testResult manager.TestResult
	testResult.Set("TestSuiteLinuxDevicesFull", device, errout, result)
	return testResult.Marshal()
}
