//+build v0.1.1
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

package specmount

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

var TestSuiteMount manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Mount"}

// TODO
func init() {
	TestSuiteMount.AddTestCase("TestMountTmpfs", TestMountTmpfs)
	TestSuiteMount.AddTestCase("TestMountDev", TestMountDev)
	TestSuiteMount.AddTestCase("TestMountCgroup", TestMountCgroup)
	TestSuiteMount.AddTestCase("TestMountDevpts", TestMountDevpts)
	TestSuiteMount.AddTestCase("TestMountMqueue", TestMountMqueue)
	TestSuiteMount.AddTestCase("TestMountShm", TestMountShm)
	TestSuiteMount.AddTestCase("TestMountSysfs", TestMountSysfs)
	manager.Manager.AddTestSuite(TestSuiteMount)
}

func setMount(fsName string, fsType string, fsSrc string, fsDes string, fsOpt []string) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	var linuxSpec specs.LinuxSpec = specsinit.SetLinuxspecMinimum()
	var linuxRuntimeSpec specs.LinuxRuntimeSpec = specsinit.SetLinuxruntimeMinimum()
	configMountTest := specs.MountPoint{fsName, fsDes}
	runtimeMountTest := specs.Mount{fsType, fsSrc, fsOpt}
	linuxSpec.Mounts = append(linuxSpec.Mounts, configMountTest)
	linuxRuntimeSpec.Mounts[fsName] = runtimeMountTest
	return linuxSpec, linuxRuntimeSpec
}

func checkHostSupport(fsname string) (bool, error) {
	cmd := exec.Command("/bin/sh", "-c", "cat /proc/filesystems | awk '{print $2}' | grep -w "+fsname)
	output, err := cmd.Output()
	if err != nil {
		return false, errors.New(err.Error())
	} else if strings.EqualFold(strings.TrimSpace(string(output)), fsname) {
		return true, nil
	} else {
		return false, nil
	}
	return true, nil
}

func testMount(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec, failinfo string) (string, error) {
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	linuxSpec.Spec.Process.Args[0] = "/bin/mount"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	out, err := adaptor.StartRunc(configFile, runtimeFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New(string(out) + err.Error())
	} else if strings.Contains(out, "/mountTest") {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New("test failed because" + failinfo)
	}
}
