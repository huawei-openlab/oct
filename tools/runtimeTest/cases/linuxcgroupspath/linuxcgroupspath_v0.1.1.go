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

package linuxcgroupspath

import (
	"errors"
	"fmt"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	// "github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var TestLinuxCgroupspath manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Cgroupspath"}

func init() {
	TestLinuxCgroupspath.AddTestCase("TestCgroupspath", TestCgroupspath)
	manager.Manager.AddTestSuite(TestLinuxCgroupspath)
}

func setCgroupspath(path string) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxRuntimeSpec := specsinit.SetLinuxruntimeMinimum()
	linuxRuntimeSpec.Linux.CgroupsPath = path
	// temporary add cgroup filesystem for test
	configMountTest := specs.MountPoint{"cgroup", "/sys/fs/cgroup"}
	runtimeMountTest := specs.Mount{"cgroup", "cgroup", []string{""}}
	linuxSpec.Mounts = append(linuxSpec.Mounts, configMountTest)
	linuxRuntimeSpec.Mounts["cgroup"] = runtimeMountTest

	return linuxSpec, linuxRuntimeSpec
}

func testCgroupspath(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec) (string, error) {
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	// check whether the container mounts Cgroup filesystem and get the mount point
	hasCgFs := false
	var cgmnt string
	for k, v := range linuxRuntimeSpec.RuntimeSpec.Mounts {
		if strings.EqualFold(v.Type, "cgroup") {
			hasCgFs = true
			for _, u := range linuxSpec.Spec.Mounts {
				if strings.EqualFold(u.Name, k) {
					cgmnt = u.Path
				}
			}
		}
	}
	if hasCgFs == false {
		return manager.UNSPPORTED, errors.New("Container doesn't support cgroup")
	}
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "find " + cgmnt + "/ -name " + linuxRuntimeSpec.Linux.CgroupsPath}
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	out, err := adaptor.StartRunc(configFile, runtimeFile)
	if err != nil {
		return manager.UNKNOWNERR, errors.New("StartRunc error :" + out + "," + err.Error())
	} else if strings.Contains(out, linuxRuntimeSpec.Linux.CgroupsPath) {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New("may be NOT SUPPORT setting cgrouppath")
	}
}

func cleanCgroup() {
	var cmd *exec.Cmd
	time.Sleep(time.Second * 15)
	cmd = exec.Command("rmdir", "/sys/fs/cgroup/*/user/*/*/runtimeValidator")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	outPut, err := cmd.Output()
	if err != nil {
		fmt.Println(string(outPut))
		log.Fatalf("[runtimeValidator] linux resources test : clean cgroup error , %v", err)
	}
	fmt.Println("clean cgroup sucess, ")
}
