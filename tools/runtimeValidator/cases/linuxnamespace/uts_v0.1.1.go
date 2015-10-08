// +build v0.1.1

// Copyright 2015 Huawei Inc. All Rights Reserved.
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

package linuxnamespace

//import "github.com/opencontainers/specs"
//shoud use github path ,but the whole project is not supported
import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/opencontainers/specs"
)

func setUtsSpec(ns specs.Namespace) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	spec := linuxSpec
	runtimespec := linuxRuntimeSpec
	runtimespec.Linux.Namespaces = append(runtimespec.Linux.Namespaces, ns)
	spec.Process.Args = []string{"/bin/readlink", "/proc/self/ns/uts"}
	return spec, runtimespec

}

func TestUtsPathEmpty() string {

	ns := specs.Namespace{Type: "uts",
		Path: ""}

	ls, lrs := setUtsSpec(ns)
	result, err := TestPathEmpty(&ls, &lrs, "/proc/*/ns/uts")

	var testResult manager.TestResult
	testResult.Set("TestUtsPathEmpty", ns, err, result)
	return testResult.Marshal()
}

func TestUtsPathUnempty() string {

	ns := specs.Namespace{Type: "uts",
		Path: "/proc/1/ns/uts"}

	ls, lrs := setUtsSpec(ns)
	result, err := TestPathUnEmpty(&ls, &lrs, ns.Path)

	var testResult manager.TestResult
	testResult.Set("TestUtsPathUnempty", ns, err, result)
	return testResult.Marshal()

}
