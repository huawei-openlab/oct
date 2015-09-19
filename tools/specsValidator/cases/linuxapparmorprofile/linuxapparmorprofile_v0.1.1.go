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

package linuxapparmorprofile

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/specsValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/specsValidator/manager"
	"github.com/huawei-openlab/oct/tools/specsValidator/utils/configconvert"
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
		Version: "0.1.0",
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
		Mounts: []specs.MountPoint{
			{
				Name: "proc",
				Path: "/proc",
			},
		},
	},
}

var linuxRuntimeSpec specs.LinuxRuntimeSpec = specs.LinuxRuntimeSpec{
	RuntimeSpec: specs.RuntimeSpec{
		Mounts: map[string]specs.Mount{
			"proc": specs.Mount{
				Type:    "proc",
				Source:  "proc",
				Options: []string{""},
			},
		},
	},
	Linux: specs.LinuxRuntime{
		Resources: &specs.Resources{
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

var TestSuiteLinuxApparmorProfile manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.ApparmorProfile"}

func init() {
	TestSuiteLinuxApparmorProfile.AddTestCase("TestLinuxApparmorProfile", TestLinuxApparmorProfile)
	manager.Manager.AddTestSuite(TestSuiteLinuxApparmorProfile)
}

func setApparmorProfile(profilename string) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxRuntimeSpec.Linux.ApparmorProfile = profilename
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "echo \" enter the container \""}
	return linuxSpec, linuxRuntimeSpec
}

func testApparmorProfile(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec) (string, error) {
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	inputprofilename := linuxRuntimeSpec.Linux.ApparmorProfile
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	cmd := exec.Command("bash", "-c", "cp cases/linuxapparmorprofile/"+inputprofilename+"  /etc/apparmor.d")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("[specsValidator] linux apparmor profile test : copy the apparmor file error, %v", err)
	}
	out, err := adaptor.StartRunc(configFile, runtimeFile)
	if err != nil {
		return manager.UNKNOWNERR, errors.New(out + err.Error())
	}
	cmd = exec.Command("bash", "-c", " apparmor_status")
	cmdout, err := cmd.Output()
	if err != nil {
		return manager.UNSPPORTED, errors.New("HOST Machine NOT Support Apparmor")
	} else if strings.Contains(strings.TrimSpace(string(cmdout)), inputprofilename) {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, nil
	}
}
