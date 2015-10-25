// +build predraft

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

package linuxuidgidmappings

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/opencontainers/specs"
	"log"
	"os"
	"os/user"
	"strconv"
)

func TestSuiteLinuxUidMappings() string {
	addTestUser()
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "cat /proc/1/uid_map"}
	//get uid&gid of test account
	testuser, _ := user.Lookup("uidgidtest")
	testuidInt, _ := strconv.ParseInt(testuser.Uid, 10, 32)
	testgidInt, _ := strconv.ParseInt(testuser.Uid, 10, 32)
	//change owner of rootfs
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatalf("utils.setBind error GOPATH == nil")
	}
	rootfspath := gopath + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/rootfs"
	utils.SetRight(rootfspath, int32(testuidInt), int32(testgidInt))
	var uid specs.IDMapping = specs.IDMapping{
		HostID:      int32(testuidInt),
		ContainerID: 0,
		Size:        10,
	}
	var gid specs.IDMapping = specs.IDMapping{
		HostID:      int32(testgidInt),
		ContainerID: 0,
		Size:        10,
	}
	failinfo := "mapping from Host UID to Container UID failed"
	linuxSpec = setIDmappings(uid, gid)
	result, err := testIDmappings(&linuxSpec, true, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestSuiteLinuxUidMappings", uid, err, result)
	cleanTestUser()
	return testResult.Marshal()

}

func TestSuiteLinuxGidMappings() string {
	addTestUser()
	linuxSpec.Spec.Process.Args = []string{"/bin/bash", "-c", "cat /proc/1/gid_map"}
	//get uid&gid of test account
	testuser, _ := user.Lookup("uidgidtest")
	testuidInt, _ := strconv.ParseInt(testuser.Uid, 10, 32)
	testgidInt, _ := strconv.ParseInt(testuser.Uid, 10, 32)
	//change owner of rootfs
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatalf("utils.setBind error GOPATH == nil")
	}
	rootfspath := gopath + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/rootfs"
	utils.SetRight(rootfspath, int32(testuidInt), int32(testgidInt))
	var uid specs.IDMapping = specs.IDMapping{
		HostID:      int32(testuidInt),
		ContainerID: 0,
		Size:        10,
	}
	var gid specs.IDMapping = specs.IDMapping{
		HostID:      int32(testgidInt),
		ContainerID: 0,
		Size:        10,
	}
	failinfo := "mapping from Host GID to Container GID failed"
	linuxSpec = setIDmappings(uid, gid)
	result, err := testIDmappings(&linuxSpec, true, failinfo)
	var testResult manager.TestResult
	testResult.Set("TestSuiteLinuxGidMappings", gid, err, result)
	cleanTestUser()
	return testResult.Marshal()

}
