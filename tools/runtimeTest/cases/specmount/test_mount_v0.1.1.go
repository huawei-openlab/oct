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
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
)

// fsName string, fsType string, fsSrc string, fsDes string, fsOpt []string

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

func TestMountCgroup() string {
	opts := []string{"nosuid", "noexec", "nodev", "relatime", "ro"}
	linuxspec, linuxruntimespec := setMount("cgroup", "cgroup", "cgroup", "/mountTest", opts)
	mount := linuxspec.Spec.Mounts
	failinfo := "Cgroup failed"
	result, err := testMount(&linuxspec, &linuxruntimespec, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestMountCgroup", mount, err, result)
	return testResult.Marshal()
}

func TestMountDev() string {
	opts := []string{}
	linuxspec, linuxruntimespec := setMount("dev", "tmpfs", "tmpfs", "/mountTest", opts)
	mount := linuxspec.Spec.Mounts
	failinfo := "dev failed"
	result, err := testMount(&linuxspec, &linuxruntimespec, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestMountCgroup", mount, err, result)
	return testResult.Marshal()
}

func TestMountDevpts() string {
	opts := []string{}
	linuxspec, linuxruntimespec := setMount("devpts", "devpts", "devpts", "/mountTest", opts)
	mount := linuxspec.Spec.Mounts
	failinfo := "devpts failed"
	result, err := testMount(&linuxspec, &linuxruntimespec, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestMountDevpts", mount, err, result)
	return testResult.Marshal()
}
func TestMountMqueue() string {
	opts := []string{}
	linuxspec, linuxruntimespec := setMount("mqueue", "mqueue", "mqueue", "/mountTest", opts)
	mount := linuxspec.Spec.Mounts
	failinfo := "mqueue failed"
	result, err := testMount(&linuxspec, &linuxruntimespec, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestMountMqueue", mount, err, result)
	return testResult.Marshal()
}

func TestMountShm() string {
	opts := []string{}
	linuxspec, linuxruntimespec := setMount("shm", "tmpfs", "shm", "/mountTest", opts)
	mount := linuxspec.Spec.Mounts
	failinfo := "Shm failed"
	result, err := testMount(&linuxspec, &linuxruntimespec, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestMountShm", mount, err, result)
	return testResult.Marshal()
}

func TestMountSysfs() string {
	opts := []string{}
	linuxspec, linuxruntimespec := setMount("sysfs", "sysfs", "sysfs", "/mountTest", opts)
	mount := linuxspec.Spec.Mounts
	failinfo := "Sysfs failed"
	result, err := testMount(&linuxspec, &linuxruntimespec, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestMountSysfs", mount, err, result)
	return testResult.Marshal()
}
