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

package linuxseccomp

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/opencontainers/specs"
	"log"
	"os"
)

func TestSuiteLinuxSeccompGetcwd() string {
	// copy the testbin into container
	var se specs.Seccomp = specs.Seccomp{
		DefaultAction: "SCMP_ACT_ALLOW",
		Syscalls: []*specs.Syscall{
			{
				Name:   "getcwd",
				Action: "SCMP_ACT_ERRNO",
			},
		},
	}
	linuxSpec = setSeccomp(se)

	result := os.Getenv("GOPATH")
	if result == "" {
		log.Fatalf("utils.setBind error GOPATH == nil")
	}
	resource := result + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/containerend"
	utils.SetRight(resource, linuxSpec.Process.User.UID, linuxSpec.Process.User.GID)
	linuxSpec.Spec.Process.Args = []string{"/testtool/linuxseccomp"}
	testtoolfolder := specs.Mount{"bind", resource, "/testtool", "bind"}
	linuxSpec.Spec.Mounts = append(linuxSpec.Spec.Mounts, testtoolfolder)
	info := ",Name=" + se.Syscalls[0].Name + ", Action=" + string(se.Syscalls[0].Action)
	result, errout := testSeccomp(&linuxSpec, info)
	var testResult manager.TestResult
	testResult.Set("TestSuiteLinuxSeccompGetcwd", se, errout, result)
	return testResult.Marshal()
}
