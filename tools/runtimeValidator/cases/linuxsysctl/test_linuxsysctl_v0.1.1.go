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

package linuxsysctl

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
)

func TestSysctlNetIpv4IpForward() string {
	sysctlInput := make(map[string]string)
	sysctlInput["net.ipv4.ip_forward"] = "4"
	failinfo := " set net.ipv4.ip_forward failed"
	var testResult manager.TestResult
	linuxspec, linuxruntimespec := setSysctls(sysctlInput)
	result, err := testSysctls(&linuxspec, &linuxruntimespec, failinfo)
	testResult.Set("TestSysctlNetIpv4IpForward", sysctlInput, err, result)
	return testResult.Marshal()
}

func TestSysctlNetCoreSomaxconn() string {
	sysctlInput := make(map[string]string)
	sysctlInput["net.core.somaxconn"] = "192"
	failinfo := " set net.core.somaxconn"
	var testResult manager.TestResult
	linuxspec, linuxruntimespec := setSysctls(sysctlInput)
	result, err := testSysctls(&linuxspec, &linuxruntimespec, failinfo)
	testResult.Set("TestSysctlNetCoreSomaxconn", sysctlInput, err, result)
	return testResult.Marshal()
}
