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

package linuxcapabilities

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
)

var TestSuiteLinuxCapabilities manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.Capabilities"}

// TODO
func init() {
	TestSuiteLinuxCapabilities.AddTestCase("TestLinuxCapabilitiesSETFCAP", TestLinuxCapabilitiesSETFCAP)
	manager.Manager.AddTestSuite(TestSuiteLinuxCapabilities)
}

func setCapability(capabilityname string) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxRuntimeSpec := specsinit.SetLinuxruntimeMinimum()
	linuxSpec.Linux.Capabilities = []string{capabilityname}
	utils.SetBind(&linuxRuntimeSpec, &linuxSpec)
	return linuxSpec, linuxRuntimeSpec
}
