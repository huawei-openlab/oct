//+build v0.1.1

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

package specmount

import (
	"github.com/huawei-openlab/oct/tools/specsValidator/manager"
)

func TestMountTmpfs() string {
	opts := []string{""}
	linuxspec, linuxruntimespec := setMount("tmpfs", "tmpfs", "tmpfs", "/mountTest", opts)
	mount := linuxspec.Spec.Mounts
	failinfo := "tmpfs failed"
	result, err := testMount(&linuxspec, &linuxruntimespec, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestMountTmpfs", mount, err, result)
	return testResult.Marshal()
}
