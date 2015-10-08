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

func setNetSpec(ns specs.Namespace) specs.LinuxSpec {
	spec := linuxSpec
	spec.Linux.Namespaces = append(spec.Linux.Namespaces, ns)
	spec.Process.Args = append(spec.Process.Args, "/bin/readlink", "/proc/self/ns/net")
	return spec

}

func TestNetPathEmpty() string {

	ns := specs.Namespace{Type: "network",
		Path: ""}

	ls := setNetSpec(ns)
	result, err := TestPathEmpty(&ls, "/proc/*/ns/net")

	var testResult manager.TestResult
	testResult.Set("TestNetPathEmpty", ns, err, result)
	return testResult.Marshal()
}

func TestNetPathUnempty() string {

	ns := specs.Namespace{Type: "network",
		Path: "/proc/1/ns/net"}

	ls := setNetSpec(ns)
	result, err := TestPathUnEmpty(&ls, ns.Path)

	var testResult manager.TestResult
	testResult.Set("TestNetPathUnempty", ns, err, result)
	return testResult.Marshal()

}
