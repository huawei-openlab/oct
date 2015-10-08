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

package linuxrootfspropagation

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
)

func TestRootfsPropagationPrivate() string {
	mkdir()
	mode := "private"
	linuxspec, linuxruntimespec := setRootfsPropagation(mode)
	result, err := testRootfsPropagationHost(&linuxspec, &linuxruntimespec, "linuxrootfspropagation")
	var testResult manager.TestResult
	testResult.Set("TestRootfsPropagationPrivate", linuxruntimespec.Linux.RootfsPropagation, err, result)
	rmdir()
	return testResult.Marshal()

}

func TestRootfsPropagationSlave() string {
	mkdir()
	mode := "slave"
	linuxspec, linuxruntimespec := setRootfsPropagation(mode)
	result, err := testRootfsPropagationHost(&linuxspec, &linuxruntimespec, "linuxrootfspropagation")
	var testResult manager.TestResult
	testResult.Set("TestRootfsPropagationSlave", linuxruntimespec.Linux.RootfsPropagation, err, result)
	rmdir()
	return testResult.Marshal()

}

func TestRootfsPropagationShare() string {
	mkdir()
	mode := "share"
	linuxspec, linuxruntimespec := setRootfsPropagation(mode)
	result, err := testRootfsPropagationHost(&linuxspec, &linuxruntimespec, "linuxrootfspropagation")
	var testResult manager.TestResult
	testResult.Set("TestRootfsPropagationShare", linuxruntimespec.Linux.RootfsPropagation, err, result)
	rmdir()
	return testResult.Marshal()
}
