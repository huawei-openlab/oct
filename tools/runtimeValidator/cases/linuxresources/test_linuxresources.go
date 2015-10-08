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

package linuxresources

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/opencontainers/specs"
	"time"
)

func TestMemoryLimit() string {
	var testResourceseMemory specs.Resources = specs.Resources{
		Memory: specs.Memory{
			Limit:       204800,
			Reservation: 0,
			Swap:        0,
			Kernel:      0,
			Swappiness:  -1,
		},
	}
	linuxspec := setResources(testResourceseMemory)
	failinfo := "Memory Limit"
	go testResources(&linuxspec)
	time.Sleep(time.Second * 3)
	defer adaptor.CleanRunc()
	result, err := checkConfigurationFromHost("memory.limit_in_bytes", "204800", failinfo)
	var testResult manager.TestResult
	testResult.Set("TestMemoryLimit", testResourceseMemory.Memory, err, result)
	return testResult.Marshal()
}
