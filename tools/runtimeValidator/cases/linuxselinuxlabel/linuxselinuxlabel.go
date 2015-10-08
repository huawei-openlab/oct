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

package linuxselinuxlabel

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/opencontainers/specs"
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

var TestSuiteLinuxSelinuxLabel manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.SelinuxProcessLabel"}

func init() {
	TestSuiteLinuxSelinuxLabel.AddTestCase("TestLinuxSelinuxProcessLabel", TestLinuxSelinuxProcessLabel)
	manager.Manager.AddTestSuite(TestSuiteLinuxSelinuxLabel)
}

func setSElinuxLabel(label string) specs.LinuxSpec {
	linuxSpec.Linux.SelinuxProcessLabel = label
	return linuxSpec
}

func testSElinuxLabel(linuxSpec *specs.LinuxSpec) (string, error) {
	selinuxlable := linuxSpec.Linux.SelinuxProcessLabel
	//test whether the host supports the selinux
	cmdout, err := exec.Command("/bin/bash", "-c", "getenforce").Output()
	if err != nil || strings.EqualFold(strings.TrimSpace(string(cmdout)), "Permissive") {
		return manager.UNSPPORTED, errors.New("Host Machine doesn't support SElinux")
	}
	configFile := "./config.json"
	linuxSpec.Process.Args = []string{"/bin/bash", "-c", "ps xZ|awk '{print $1}' |sed -n '2p' "}
	err = configconvert.LinuxSpecToConfig(configFile, linuxSpec)

	out, err := adaptor.StartRunc(configFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New("StartRunc error :" + out + "," + err.Error())
	} else {
		if strings.EqualFold(strings.TrimSpace(string(out)), label) {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New("test failed because selinux label is not effective")
		}
	}
}
