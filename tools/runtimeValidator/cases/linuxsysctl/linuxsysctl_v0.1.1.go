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

package linuxsysctl

import (
	"errors"
	"fmt"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"strings"
)

var TestSuiteLinuxSysctl manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Sysctl"}

func init() {
	TestSuiteLinuxSysctl.AddTestCase("TestSysctlNetIpv4IpForward", TestSysctlNetIpv4IpForward)
	TestSuiteLinuxSysctl.AddTestCase("TestSysctlNetCoreSomaxconn", TestSysctlNetCoreSomaxconn)
}

func setSysctls(testsysctls map[string]string) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxRuntimeSpec := specsinit.SetLinuxruntimeMinimum()
	linuxRuntimeSpec.Linux.Sysctl = testsysctls
	fmt.Println("setsysctls done")
	return linuxSpec, linuxRuntimeSpec
}

func testSysctls(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec, failinfo string) (string, error) {
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	var key, value string
	for k, v := range linuxRuntimeSpec.Linux.Sysctl {
		linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "sysctl " + k}
		key = k
		value = v
	}
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	out, err := adaptor.StartRunc(configFile, runtimeFile)
	fmt.Println("out=" + out)
	if err != nil {
		return manager.UNSPPORTED, errors.New(string(out) + err.Error())
	} else {
		if strings.EqualFold(strings.TrimSpace(out), key+" = "+value) {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New("test failed because" + failinfo)
		}
	}
}
