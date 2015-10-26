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

package linuxrlimits

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"strings"
)

var TestSuiteLinuxRlimits manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Rlimits"}

func init() {
	TestSuiteLinuxRlimits.AddTestCase("TestRlimitNOFILESoft", TestRlimitNOFILESoft)
	TestSuiteLinuxRlimits.AddTestCase("TestRlimitNOFILEHard", TestRlimitNOFILEHard)
	manager.Manager.AddTestSuite(TestSuiteLinuxRlimits)
}

func setRlimits(testrlimits specs.Rlimit) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxRuntimeSpec := specsinit.SetLinuxruntimeMinimum()
	linuxRuntimeSpec.Linux.Rlimits = []specs.Rlimit{testrlimits}
	return linuxSpec, linuxRuntimeSpec
}

func testRlimits(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec, rlimitItem string, value string, isSoftLimit bool) (string, error) {
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	if isSoftLimit {
		linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "ulimit " + rlimitItem + " -S"}
	} else {
		linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "ulimit " + rlimitItem + " -H"}
	}
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	out, err := adaptor.StartRunc(configFile, runtimeFile)
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
