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

package linuxrootfspropagation

import (
	"github.com/huawei-openlab/oct/tools/specsValidator/manager"
)

func TestRootfsPropagationPrivate() string {
	mode := "private"
	linuxspec := setRootfsPropagation(mode)
	result, err := testRootfsPropagationHost(&linuxspec, "guest_linuxrootfspropagation")
	var testResult manager.TestResult
	testResult.Set("TestRootfsPropagationPrivate", linuxspec.Linux.RootfsPropagation, err, result)
	return testResult.Marshal()
}

func TestRootfsPropagationSlave() string {
	mode := "slave"
	linuxspec := setRootfsPropagation(mode)
	result, err := testRootfsPropagationHost(&linuxspec, "guest_linuxrootfspropagation")
	var testResult manager.TestResult
	testResult.Set("TestRootfsPropagationSlave", linuxspec.Linux.RootfsPropagation, err, result)
	return testResult.Marshal()
}

func TestRootfsPropagationShare() string {
	mode := "share"
	linuxspec := setRootfsPropagation(mode)
	result, err := testRootfsPropagationHost(&linuxspec, "guest_linuxrootfspropagation")
	var testResult manager.TestResult
	testResult.Set("TestRootfsPropagationShare", linuxspec.Linux.RootfsPropagation, err, result)
	return testResult.Marshal()
}
