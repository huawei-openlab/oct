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

package linuxresources

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/opencontainers/specs"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
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

var TestSuiteLinuxResources manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Resources"}

func init() {
	TestSuiteLinuxResources.AddTestCase("TestMemoryLimit", TestMemoryLimit)
	manager.Manager.AddTestSuite(TestSuiteLinuxResources)
}

func setResources(resources specs.Resources) specs.LinuxSpec {
	linuxSpec.Linux.Resources = resources
	return linuxSpec
}

func testResources(linuxSpec *specs.LinuxSpec) (string, error) {
	fmt.Println("enter test source")
	configFile := "./config.json"
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "sleep 10s"}
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	out, err := adaptor.StartRunc(configFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New("StartRunc error :" + out + "," + err.Error())
	} else {
		return manager.PASSED, nil
	}
}

func checkConfigurationFromHost(filename string, configvalue string, failinfo string) (string, error) {
	cmd := exec.Command("bash", "-c", "cat  /sys/fs/cgroup/*/*/*/*/runtimeValidator/"+filename)
	cmdouput, err := cmd.Output()
	if err != nil {
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

func cleanCgroup() {
	// cmd := exec.Command("bash", "-c", "rmdir /sys/fs/cgroup/memory/user/1002.user/c2.session/runtimeValidator")
	// outPut, err := cmd.Output()
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
