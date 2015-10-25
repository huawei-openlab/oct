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

package linuxuidgidmappings

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/opencontainers/specs"
	"log"
	"os/exec"
	// "os/user"
	"runtime"
	"strconv"
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
			{
				Type: "user",
				Path: "",
			},
			{
				Type: "pid",
				Path: "",
			},
		},
	},
}

var TestSuiteLinuxUidGidMappings manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.UidGidMappings"}

// TODO
func init() {
	TestSuiteLinuxUidGidMappings.AddTestCase("TestSuiteLinuxUidMappings", TestSuiteLinuxUidMappings)
	TestSuiteLinuxUidGidMappings.AddTestCase("TestSuiteLinuxGidMappings", TestSuiteLinuxGidMappings)
	manager.Manager.AddTestSuite(TestSuiteLinuxUidGidMappings)
}

func setIDmappings(testuid specs.IDMapping, testgid specs.IDMapping) specs.LinuxSpec {
	linuxSpec.Linux.UIDMappings = []specs.IDMapping{testuid}
	linuxSpec.Linux.GIDMappings = []specs.IDMapping{testgid}
	return linuxSpec
}

func testIDmappings(linuxSpec *specs.LinuxSpec, isUid bool, failinfo string) (string, error) {
	//test whether the usernamespace works
	configFile := "./config.json"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	out, err := adaptor.StartRunc(configFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New("StartRunc error :" + out + "," + err.Error())
	}
	outarray := strings.Fields(strings.TrimSpace(out))
	outcuid, _ := strconv.ParseInt(outarray[0], 10, 0)
	outhuid, _ := strconv.ParseInt(outarray[1], 10, 0)
	outsize, _ := strconv.ParseInt(outarray[2], 10, 0)
	var incuid, inhuid, insize int32
	if isUid {
		incuid = linuxSpec.Linux.UIDMappings[0].ContainerID
		inhuid = linuxSpec.Linux.UIDMappings[0].HostID
		insize = linuxSpec.Linux.UIDMappings[0].Size
	} else {
		incuid = linuxSpec.Linux.GIDMappings[0].ContainerID
		inhuid = linuxSpec.Linux.GIDMappings[0].HostID
		insize = linuxSpec.Linux.GIDMappings[0].Size
	}
	if (int32(outcuid) == incuid) && (int32(outhuid) == inhuid) && (int32(outsize) == insize) {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New("test failed because" + failinfo)
	}
}
func addTestUser() {
	cmd := exec.Command("/bin/bash", "-c", "useradd -m uidgidtest")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux linux uid/gid mappings test : create test user account error, %v", err)
	}
}

func cleanTestUser() {
	cmd := exec.Command("/bin/bash", "-c", "userdel -r uidgidtest")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux linux uid/gid mappings test : delete test user account error, %v", err)
	}
}
