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

package linuxresources

import (
	"errors"
	"fmt"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var TestSuiteLinuxResources manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Resources"}

func init() {
	TestSuiteLinuxResources.AddTestCase("TestResourceCpuQuota", TestCpuQuota)
	TestSuiteLinuxResources.AddTestCase("TestResourceBlockIOWeight", TestBlockIOWeight)
	TestSuiteLinuxResources.AddTestCase("TestResourceMemoryLimit", TestMemoryLimit)
	TestSuiteLinuxResources.AddTestCase("TestResourceHugepageLimit", TestHugepageLimit)
	manager.Manager.AddTestSuite(TestSuiteLinuxResources)
}

func setResources(resources specs.Resources) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxRuntimeSpec := specsinit.SetLinuxruntimeMinimum()
	linuxRuntimeSpec.Linux.Resources = &resources
	return linuxSpec, linuxRuntimeSpec
}

func testResources(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec) (string, error) {
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "sleep 0.5s"}
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	out, err := adaptor.StartRunc(configFile, runtimeFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New("StartRunc error :" + out + "," + err.Error())
	} else {
		return manager.PASSED, nil
	}
}

func checkConfigurationFromHost(subsys string, filename string, configvalue string, failinfo string) (string, error) {
	time.Sleep(time.Millisecond * 100)
	pwd, err := exec.Command("bash", "-c", "find /sys/fs/cgroup/"+subsys+" -name runtimeValidator").Output()
	cmdouput, err := exec.Command("bash", "-c", "cat "+strings.TrimSpace(string(pwd))+"/"+filename).Output()
	if err != nil {
		log.Println("path=" + strings.TrimSpace(string(pwd)) + "/" + filename)
		log.Fatalf("[runtimeValidator] linux resources test : read the "+filename+" error, %v", err)
		return manager.UNKNOWNERR, err
	} else {
		if strings.EqualFold(strings.TrimSpace(string(cmdouput)), configvalue) {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New("test failed because" + failinfo)
		}
	}

}

func cleanCgroup(path string) {
	utils.ExecCmdStr("apt-get", "/", "install", "-y", "cgroup-bin")
	var cmd *exec.Cmd
	cmd = exec.Command("cgdelete", path)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	outPut, err := cmd.Output()
	if err != nil {
		fmt.Println(string(outPut))
		log.Fatalf("[runtimeValidator] linux resources test : clean cgroup error , %v", err)
	}
}
