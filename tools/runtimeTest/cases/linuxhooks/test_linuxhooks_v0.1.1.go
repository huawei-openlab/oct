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

package linuxhooks

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/opencontainers/specs"
)

func TestLinuxHooksPrestart() string {
	var pre []specs.Hook = []specs.Hook{
		{
			Path: "/bin/bash",
			Args: []string{"-c", "touch /testHooksPrestart.txt"},
			Env:  []string{""},
		},
	}
	var testResult manager.TestResult
	linuxspec, linuxruntimespec := setHooks(pre, true)
	info := "may be not support"
	result, err := testHooks(&linuxspec, &linuxruntimespec, "testHooksPrestart.txt", info)
	testResult.Set("TestSuiteLinuxHooksPrestart", pre, err, result)
	return testResult.Marshal()
}
