// +build v0.1.1

// Copyright 2014 Google Inc. All Rights Reserved.
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

package specversion

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
)

// The test func for TestCase TestVersionCorrect
func TestVersionCorrect() string {

	// Set result to spec.Version, and get specs.LinuxSpec obj
	ls := setVersion(testValueCorrect)

	// Get smallest specs.LinuxRuntimeSpec obj
	lr := specsinit.SetLinuxruntimeMinimum()
	version := ls.Spec.Version

	// Do test
	result, err := testVersion(&ls, &lr, true)
	var testResult manager.TestResult

	// Set reusult to TestResult
	testResult.Set("TestVersionCorrect", version, err, result)
	return testResult.Marshal()
}
