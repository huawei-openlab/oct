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

package linuxseccomp

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/opencontainers/specs"
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
	linuxspec, linuxruntimespec := setSeccomp(se)

	utils.SetBind(&linuxruntimespec, &linuxspec)
	linuxspec.Spec.Process.Args = []string{"/bin/bash", "-c", "/containerend/linuxseccomp"}
	info := ",Name=" + se.Syscalls[0].Name + ", Action=" + string(se.Syscalls[0].Action)
	result, errout := testSeccomp(&linuxspec, &linuxruntimespec, info)
	var testResult manager.TestResult
	testResult.Set("TestSuiteLinuxSeccompGetcwd", se, errout, result)
	return testResult.Marshal()
}
