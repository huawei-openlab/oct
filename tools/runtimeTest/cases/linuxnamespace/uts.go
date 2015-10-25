// +build predraft

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

func setUtsSpec(ns specs.Namespace) specs.LinuxSpec {
	spec := linuxSpec
	spec.Linux.Namespaces = append(spec.Linux.Namespaces, ns)
	spec.Process.Args = append(spec.Process.Args, "/bin/readlink", "/proc/self/ns/uts")
	return spec

}

func TestUtsPathEmpty() string {

	ns := specs.Namespace{Type: "uts",
		Path: ""}

	ls := setUtsSpec(ns)
	result, err := TestPathEmpty(&ls, "/proc/*/ns/uts")

	var testResult manager.TestResult
	testResult.Set("TestUtsPathEmpty", ns, err, result)
	return testResult.Marshal()
}

func TestUtsPathUnempty() string {

	ns := specs.Namespace{Type: "uts",
		Path: "/proc/1/ns/uts"}

	ls := setUtsSpec(ns)
	result, err := TestPathUnEmpty(&ls, ns.Path)

	var testResult manager.TestResult
	testResult.Set("TestUtsPathUnempty", ns, err, result)
	return testResult.Marshal()

}
