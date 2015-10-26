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

package linuxcapabilities

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"strings"
)

func TestLinuxCapabilitiesSETFCAP() string {
	linuxspec := setCapability("SETFCAP")
	linuxspec.Spec.Process.Args = []string{"/sbin/setcap", "CAP_SETFCAP=eip", "/testtool/linuxcapabilities"}
	capability := linuxspec.Linux.Capabilities
	configFile := "./config.json"
	err := configconvert.LinuxSpecToConfig(configFile, &linuxspec)
	out, err := adaptor.StartRunc(configFile)
	var result string
	var errout error
	if err != nil {
		result = manager.UNSPPORTED
		errout = errors.New(string(out) + err.Error())
	} else if strings.EqualFold(strings.TrimSpace(string(out)), "") {
		result = manager.PASSED
		errout = nil
	} else {
		result = manager.FAILED
		errout = errors.New("test Capabilities CAP_SETFCAP NOT  does''t work")
	}
	var testResult manager.TestResult
	testResult.Set("TestMountTmpfs", capability, errout, result)
	return testResult.Marshal()
}
