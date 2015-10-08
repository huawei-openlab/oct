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

package specmount

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

var TestSuiteMount manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Mount"}

// TODO
func init() {
	TestSuiteMount.AddTestCase("TestMountTmpfs", TestMountTmpfs)
	manager.Manager.AddTestSuite(TestSuiteMount)
}

func setMount(fsName string, fsSrc string, fsDes string, fsOpt string) specs.LinuxSpec {
	mountsorigin := specs.Mount{"proc", "proc", "/proc", ""}
	mountsadd := specs.Mount{fsName, fsSrc, fsDes, fsOpt}
	linuxSpec.Mounts = []specs.Mount{mountsorigin, mountsadd}
	return linuxSpec
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

func testMount(linuxSpec *specs.LinuxSpec, failinfo string) (string, error) {
	configFile := "./config.json"
	linuxSpec.Spec.Process.Args[0] = "/bin/mount"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	out, err := adaptor.StartRunc(configFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New(string(out) + err.Error())
	} else if strings.Contains(out, "/mountTest") {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New("test failed because" + failinfo)
	}
}
