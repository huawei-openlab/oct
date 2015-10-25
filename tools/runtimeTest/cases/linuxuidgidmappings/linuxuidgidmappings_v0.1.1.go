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

package linuxuidgidmappings

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var TestLinuxUidGidMappings manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.UidGidMappings"}

// TODO
func init() {
	TestLinuxUidGidMappings.AddTestCase("TestLinuxUidMappings", TestLinuxUidMappings)
	TestLinuxUidGidMappings.AddTestCase("TestLinuxGidMappings", TestLinuxGidMappings)
	manager.Manager.AddTestSuite(TestLinuxUidGidMappings)
}

func setIDmappings(testuid specs.IDMapping, testgid specs.IDMapping) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxRuntimeSpec := specsinit.SetLinuxruntimeMinimum()
	userns := specs.Namespace{specs.UserNamespace, ""}
	pidns := specs.Namespace{specs.PIDNamespace, ""}
	linuxRuntimeSpec.Linux.Namespaces = append(linuxRuntimeSpec.Linux.Namespaces, userns, pidns)
	linuxRuntimeSpec.Linux.UIDMappings = []specs.IDMapping{testuid}
	linuxRuntimeSpec.Linux.GIDMappings = []specs.IDMapping{testgid}
	return linuxSpec, linuxRuntimeSpec
}

func testIDmappings(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec, isUid bool, failinfo string) (string, error) {
	//test whether the usernamespace works
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	out, err := adaptor.StartRunc(configFile, runtimeFile)
	if err != nil {
		return manager.UNSPPORTED, errors.New("StartRunc error :" + out + "," + err.Error())
	}
	outarray := strings.Fields(strings.TrimSpace(out))
	outcuid, _ := strconv.ParseInt(outarray[0], 10, 0)
	outhuid, _ := strconv.ParseInt(outarray[1], 10, 0)
	outsize, _ := strconv.ParseInt(outarray[2], 10, 0)
	var incuid, inhuid, insize uint32
	if isUid {
		incuid = linuxRuntimeSpec.Linux.UIDMappings[0].ContainerID
		inhuid = linuxRuntimeSpec.Linux.UIDMappings[0].HostID
		insize = linuxRuntimeSpec.Linux.UIDMappings[0].Size
	} else {
		incuid = linuxRuntimeSpec.Linux.GIDMappings[0].ContainerID
		inhuid = linuxRuntimeSpec.Linux.GIDMappings[0].HostID
		insize = linuxRuntimeSpec.Linux.GIDMappings[0].Size
	}
	if (uint32(outcuid) == incuid) && (uint32(outhuid) == inhuid) && (uint32(outsize) == insize) {
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
