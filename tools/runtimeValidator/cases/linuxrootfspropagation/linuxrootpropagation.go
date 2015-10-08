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

package linuxrootfspropagation

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/opencontainers/specs"
	"log"
	"os"
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
		Devices: []specs.Device{
			{
				Type:        99,
				Path:        "/dev/null",
				Major:       1,
				Minor:       3,
				Permissions: "rwm",
				FileMode:    438,
				UID:         0,
				GID:         0,
			},
		},
	},
}

const (
	SLAVE   = "slave"
	PRIVATE = "private"
	SHARE   = "share"
)

var TestRootfsPropagation manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.RootfsPropagation"}

func init() {
	TestRootfsPropagation.AddTestCase("TestRootfsPropagationPrivate", TestRootfsPropagationPrivate)
	TestRootfsPropagation.AddTestCase("TestRootfsPropagationSlave", TestRootfsPropagationSlave)
	TestRootfsPropagation.AddTestCase("TestRootfsPropagationShare", TestRootfsPropagationShare)
	manager.Manager.AddTestSuite(TestRootfsPropagation)
}

func setRootfsPropagation(mode string) specs.LinuxSpec {
	linuxSpec.Linux.RootfsPropagation = mode
	return linuxSpec
}

func testRootfsPropagationHost(linuxSpec *specs.LinuxSpec, guestfilename string) (string, error) {
	//
	configFile := "./config.json"
	propagationmode := linuxSpec.Linux.RootfsPropagation

	cmd := exec.Command("bash", "-c", "touch  rootfs/fspropagationtest/fromhost.txt")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux rootfs propagation test : touch test file in host error, %v", err)
	}
	// set the config parameters relative to this case
	result := os.Getenv("GOPATH")
	if result == "" {
		log.Fatalf("utils.setBind error GOPATH == nil")
	}
	resource := result + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/containerend"
	utils.SetRight(resource, linuxSpec.Process.User.UID, linuxSpec.Process.User.GID)
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "/testtool/" + guestfilename}
	testtoolfolder := specs.Mount{"bind", resource, "/testtool", "bind"}
	linuxSpec.Spec.Mounts = append(linuxSpec.Spec.Mounts, testtoolfolder)
	linuxSpec.Linux.Capabilities = []string{"AUDIT_WRITE", "KILL", "NET_BIND_SERVICE", "SYS_ADMIN"}
	linuxSpec.Spec.Root.Readonly = false

	err = configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	out_container, err := adaptor.StartRunc(configFile)
	cmd = exec.Command("/bin/bash", "-c", "ls rootfs/fspropagationtest")
	out_host, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux rootfs propagation test : read test file from container (in host) error, %v", err)
		return manager.UNKNOWNERR, err
	}
	var flag_container, flag_host bool
	if strings.Contains(strings.TrimSpace(out_container), "fromhost.txt") {
		flag_container = true
	} else {
		flag_container = false
	}
	if strings.Contains(strings.TrimSpace(string(out_host)), "fromcontainer.txt") {
		flag_host = true
	} else {
		flag_container = false
	}
	switch propagationmode {
	case "slave":
		if flag_container == true && flag_host == false {
			return manager.PASSED, nil
		}
	case "private":
		if flag_container == false && flag_host == false {
			return manager.PASSED, nil
		}
	case "share":
		if flag_container && flag_host {
			return manager.PASSED, nil
		}
	}
	return manager.FAILED, errors.New("RootfsPropagationmode:" + propagationmode + "failed")
}

func mkdir() {
	cmd := exec.Command("bash", "-c", "mkdir  rootfs/fspropagationtest ")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux rootfs propagation test : make new folder in host error, %v", err)
	}
}

func rmdir() {
	cmd := exec.Command("bash", "-c", " rm -r rootfs/fspropagationtest/ ")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux rootfs propagation test : remove folder in host error, %v", err)
	}
}
