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

package linuxapparmorprofile

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
)

var linuxSpec specs.LinuxSpec = specsinit.SetLinuxspecMinimum()
var linuxRuntimeSpec specs.LinuxRuntimeSpec = specsinit.SetLinuxruntimeMinimum()

var TestSuiteLinuxApparmorProfile manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Linux.ApparmorProfile"}

func init() {
	TestSuiteLinuxApparmorProfile.AddTestCase("TestLinuxApparmorProfile", TestLinuxApparmorProfile)
	manager.Manager.AddTestSuite(TestSuiteLinuxApparmorProfile)
}

func setApparmorProfile(profilename string) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	linuxRuntimeSpec.Linux.ApparmorProfile = profilename
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "sleep 0.5s"}
	return linuxSpec, linuxRuntimeSpec
}

func testApparmorProfile(linuxSpec *specs.LinuxSpec, linuxRuntimeSpec *specs.LinuxRuntimeSpec) (string, error) {
	out, err := pretest()
	configFile := "./config.json"
	runtimeFile := "./runtime.json"
	err = configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	err = configconvert.LinuxRuntimeToConfig(runtimeFile, linuxRuntimeSpec)
	c := make(chan bool)
	go func() {
		adaptor.StartRunc(configFile, runtimeFile)
		close(c)
	}()
	out, err = checkapparmorfilefromhost()
	<-c
	if err != nil {
		return manager.UNKNOWNERR, errors.New(out + err.Error())
	}
	return out, err

}

func pretest() (string, error) {
	out, err := ioutil.ReadFile("/sys/module/apparmor/parameters/enabled")
	if err != nil || !strings.EqualFold(strings.TrimSpace(string(out)), "Y") {
		return manager.UNSPPORTED, errors.New("HOST Machine NOT Support Apparmor")
	}
	out, err = exec.Command("apparmor_parser", "-r", "cases/linuxapparmorprofile/testapporprofile ").Output()
	if err != nil {
		return manager.UNKNOWNERR, errors.New("HOST Machine Load Apparmor ERROR")
	}
	return manager.PASSED, nil
}

func checkapparmorfilefromhost() (string, error) {
	time.Sleep(time.Millisecond * 100)
	out, err := exec.Command("apparmor_status").Output()
	outstr := string(out)
	outstr = strings.TrimLeft(outstr, "processes are in enforce mode")
	outstr = strings.TrimRight(outstr, "processes are in complain mode")
	if err != nil {
		return manager.UNKNOWNERR, errors.New("check apparmor status error")
	} else {
		if strings.Contains(outstr, "runc-test") {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New("apparmorfile does not work")
		}
	}
}
