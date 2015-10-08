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

package linuxcgroupspath

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
)

func TestCgroupspath() string {
	cgpath := "/cgroupspathtest/testcontainer"
	linuxspec, linuxruntimespec := setCgroupspath(cgpath)
	result, err := testCgroupspath(&linuxspec, &linuxruntimespec)
	var testResult manager.TestResult
	testResult.Set("TestCgroupspath", cgpath, err, result)
	return testResult.Marshal()
}
