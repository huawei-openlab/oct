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

package linuxselinuxlabel

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"os/exec"
	"strings"
)

var TestSuiteLinuxSelinuxLabel manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.SelinuxProcessLabel"}

func init() {
	TestSuiteLinuxSelinuxLabel.AddTestCase("TestLinuxSelinuxProcessLabel", TestLinuxSelinuxProcessLabel)
	manager.Manager.AddTestSuite(TestSuiteLinuxSelinuxLabel)
}

func setSElinuxLabel(label string) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxRuntimeSpec := specsinit.SetLinuxruntimeMinimum()
	linuxRuntimeSpec.Linux.SelinuxProcessLabel = label
	return linuxSpec, linuxRuntimeSpec
}

func testSElinuxLabel(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec) (string, error) {
	selinuxlable := linuxRuntimeSpec.Linux.SelinuxProcessLabel
	//test whether the host supports the selinux
	cmdout, err := exec.Command("/bin/bash", "-c", "getenforce").Output()
	if err != nil || strings.EqualFold(strings.TrimSpace(string(cmdout)), "Permissive") {
		return manager.UNSPPORTED, errors.New("Host Machine doesn't support SElinux")
	}
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	linuxSpec.Process.Args = []string{"/bin/bash", "-c", "ps xZ|awk '{print $1}' |sed -n '2p' "}
	err = configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	out, err := adaptor.StartRunc(configFile, runtimeFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New("StartRunc error :" + out + "," + err.Error())
	} else {
		if strings.EqualFold(strings.TrimSpace(string(out)), selinuxlable) {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New("test failed because selinux label is not effective")
		}
	}
}
