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

package linuxcapabilities

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/opencontainers/specs"
	"log"
	"os"
	"runtime"
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

var TestSuiteLinuxCapabilities manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Capabilities"}

// TODO
func init() {
	TestSuiteLinuxCapabilities.AddTestCase("TestLinuxCapabilitiesSETFCAP", TestLinuxCapabilitiesSETFCAP)
	manager.Manager.AddTestSuite(TestSuiteLinuxCapabilities)
}

func setCapability(capabilityname string) specs.LinuxSpec {
	linuxSpec.Linux.Capabilities = []string{capabilityname}
	result := os.Getenv("GOPATH")
	if result == "" {
		log.Fatalf("utils.setBind error GOPATH == nil")
	}
	resource := result + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/containerend"
	utils.SetRight(resource, linuxSpec.Process.User.UID, linuxSpec.Process.User.GID)
	testtoolfolder := specs.Mount{"bind", resource, "/testtool", "bind"}
	linuxSpec.Spec.Mounts = append(linuxSpec.Spec.Mounts, testtoolfolder)
	linuxSpec.Process.Cwd = "/testtool"
	return linuxSpec
}
