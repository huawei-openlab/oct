// +build predraft

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

package specroot

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
)

func TestReadOnlyTrue() string {
	linuxspec := setRoot(true, testPathCorrect)
	root := linuxspec.Spec.Root
	result, err := testRoot(&linuxspec, true, testPathCorrect)
	var testResult manager.TestResult
	testResult.Set("TestReadOnlyTrue", root, err, result)
	return testResult.Marshal()
}
func TestReadOnlyFalse() string {
	linuxspec := setRoot(false, testPathCorrect)
	root := linuxspec.Spec.Root
	result, err := testRoot(&linuxspec, false, testPathCorrect)
	var testResult manager.TestResult
	testResult.Set("TestReadOnlyFalse", root, err, result)
	return testResult.Marshal()
}

func TestPathError() string {
	linuxspec := setRoot(true, testPathError)
	root := linuxspec.Spec.Root
	result, err := testRoot(&linuxspec, true, testPathError)
	var testResult manager.TestResult
	testResult.Set("TestPathError", root, err, result)
	return testResult.Marshal()
}
