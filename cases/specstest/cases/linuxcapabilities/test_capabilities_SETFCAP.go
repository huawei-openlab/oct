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

package linuxcapabilities

import (
	"errors"
	"github.com/huawei-openlab/oct/cases/specstest/adaptor"
	"github.com/huawei-openlab/oct/cases/specstest/manager"
	"github.com/huawei-openlab/oct/cases/specstest/utils/configconvert"
	"log"
	"os/exec"
	"strings"
)

func TestLinuxCapabilitiesSETFCAP() string {
	// copy the testbin into container
	cmd := exec.Command("/bin/sh", "-c", "cp  cases/linuxcapabilities/capabilitytestbin /tmp/testtool")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux Capabilities test : init the testbin file error, %v", err)
	}

	linuxspec := setCapability("SETFCAP")
	linuxspec.Spec.Process.Args = []string{"/sbin/setcap", "CAP_SETFCAP=eip", "/testtool/capabilitytestbin"}
	capability := linuxspec.Linux.Capabilities
	configFile := "./config.json"
	err = configconvert.LinuxSpecToConfig(configFile, &linuxspec)
	out, err := adaptor.StartRunc(configFile)
	var result string
	var errout error
	if err != nil {
		result = manager.UNSPPORTED
		errout = errors.New(string(out) + err.Error())
	} else if strings.EqualFold(strings.TrimSpace(string(out)), "") {
		result = manager.PASSED
		errout = nil
	} else {
		result = manager.FAILED
		errout = nil
	}
	var testResult manager.TestResult
	testResult.Set("TestMountTmpfs", capability, errout, result)
	return testResult.Marshal()
}
