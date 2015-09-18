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
	"github.com/huawei-openlab/oct/tools/specsValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/specsValidator/manager"
	"github.com/huawei-openlab/oct/tools/specsValidator/utils/configconvert"
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

var TestSuiteMount manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Mount"}

// TODO
func init() {
	TestSuiteMount.AddTestCase("TestMountTmpfs", TestMountTmpfs)
	manager.Manager.AddTestSuite(TestSuiteMount)
}

func setMount(fsName string, fsType string, fsSrc string, fsDes string, fsOpt []string) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
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
	out, err := adaptor.StartRunc("", "")
	if err != nil {
		return manager.UNSPPORTED, errors.New(string(out) + err.Error())
	} else if strings.Contains(out, "/mountTest") {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New("test failed because" + failinfo)
	}
}
