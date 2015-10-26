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

package linuxselinuxlabel

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
)

func TestLinuxSelinuxProcessLabel() string {
	labelin := "system_u:system_r:svirt_lxc_net_t:s0:c124,c675"
	linuxspec, linuxruntimespec := setSElinuxLabel(labelin)
	result, err := testSElinuxLabel(&linuxspec, &linuxruntimespec)
	var testResult manager.TestResult
	testResult.Set("TestLinuxSelinuxProcessLabel", linuxruntimespec.Linux.SelinuxProcessLabel, err, result)
	return testResult.Marshal()
}
