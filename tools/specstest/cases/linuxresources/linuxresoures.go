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
	"github.com/huawei-openlab/oct/tools/specstest/adaptor"
	"github.com/huawei-openlab/oct/tools/specstest/manager"
	"github.com/huawei-openlab/oct/tools/specstest/utils/configconvert"
	"github.com/opencontainers/specs"
	"log"
	"os/exec"
	"runtime"
	"strings"
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
			Path:     "rootfs_rootconfig",
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
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "sleep 30s"}
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	out, err := adaptor.StartRunc(configFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New(string(out) + err.Error())
	} else {
		fmt.Println("runc start success")
		return manager.PASSED, nil
	}
}

func checkConfigurationFromHost(filename string, configvalue string) (string, error) {
	cmd := exec.Command("bash", "-c", "cat  /sys/fs/*/*/*/*/*/specstest/"+filename)
	cmdouput, err := cmd.Output()
	fmt.Println("cmdoutput=" + string(cmdouput))
	if err != nil {
		log.Fatalf("[Specstest] linux resources test : read the "+filename+" error, %v", err)
		return manager.UNKNOWNERR, nil
	} else {
		if strings.EqualFold(strings.TrimSpace(string(cmdouput)), configvalue) {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, nil
		}
	}
}
