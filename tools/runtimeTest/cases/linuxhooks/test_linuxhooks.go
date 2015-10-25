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

/*
"hooks": {
    "prestart": [{
      "path": "/bin/ls",
      "args": null,
      "env": null
    }],
    "poststop": [{
      "path": "/bin/pwd",
      "args": null,
      "env": null
    }]
  },
*/

package linuxhooks

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/opencontainers/specs"
)

func TestSuiteLinuxHooksPrestart() string {
	var pre []specs.Hook = []specs.Hook{
		{
			Path: "/bin/bash",
			Args: []string{"-c", "touch /testHooksPrestart.txt"},
			Env:  []string{""},
		},
	}
	var testResult manager.TestResult
	linuxspec := setHooks(pre, true)
	info := " :Prestart Path=" + pre[0].Path
	result, err := testHooks(&linuxspec, "testHooksPrestart.txt", info)
	testResult.Set("TestSuiteLinuxHooksPrestart", linuxspec.Hooks.Prestart, err, result)
	return testResult.Marshal()
}
