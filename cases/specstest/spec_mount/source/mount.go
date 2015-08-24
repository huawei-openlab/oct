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

package main

import (
	configconvert "./../../source/configconvert"
	hostsetup "./../../source/hostsetup"
	runcstart "./../../source/runcstart"
	specs "./../../source/specs"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type TestResult struct {
	Mount MountStruct `json:"Linuxspec.Spec.Mount"`
}

// new struct to output testresult from the origin Mount struct
type MountStruct struct {
	// Type specifies the mount kind.
	Type map[string]string `json:"type"`
	// Source specifies the source path of the mount. In the case of bind mounts on
	// linux based systems this would be the file on the host.
	Source map[string]string `json:"source"`
	// Destination is the path where the mount will be placed relative to the container's root.
	Destination map[string]string `json:"destination"`
	// Options are fstab style mount options.
	Options map[string]string `json:"options"`
}

// update config.json according to the test case
func updateConfig(fsName string, fsSrc string, fsDes string, fsOpt string, configFilepath string) error {
	var linuxspec *specs.LinuxSpec
	linuxspec, err := configconvert.ConfigToLinuxSpec(configFilepath)
	if err != nil {
		log.Fatalf("[Specstestroot] Root test readonly = %v readconfig error, err = %v", fsName, err)
	}
	mountsorigin := specs.Mount{"proc", "proc", "/proc", ""}
	mountsadd := specs.Mount{fsName, fsSrc, fsDes, fsOpt}
	mountsbind := specs.Mount{"bind", "/tmp/testtool", "/testtool", "rbind,rw"}
	mountsnew := []specs.Mount{mountsorigin, mountsbind, mountsadd}
	linuxspec.Mounts = mountsnew
	linuxspec.Spec.Process.Args = []string{"./mount_guest"}
	linuxspec.Spec.Root.Path = "./../../source/rootfs_rootconfig"
	err = configconvert.LinuxSpecToConfig(configFilepath, linuxspec)
	fmt.Println(linuxspec)
	return err
}

func checkHostSupport(fsname string) bool {
	cmd := exec.Command("/bin/sh", "-c", "cat /proc/filesystems | awk '{print $2}' | grep -w "+fsname)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] Mount test filesystemtype= %v cat /proc/filesystems failed, err = %v...", fsname, err)
	} else if strings.EqualFold(strings.TrimSpace(string(output)), fsname) {
		return true
	} else {
		return false
	}
	return true
}

func mountSupportTest(fsname string, fssrc string, fsdes string, fsopt string) (string, error) {
	configFile := "./../../source/config.json"
	err := updateConfig(fsname, fssrc, fsdes, fsopt, configFile)
	var resultString string
	if err != nil {
		log.Fatalf("[Specstest] Mount test filesystemtype= %v setupEnv failed, err = %v...", fsname, err)
	} else {
		log.Printf("[Specstest] Mount test filesystemtype= %v setupEnv sucess ... \n", fsname)
	}

	output, err := runcstart.StartRunc(configFile)
	fmt.Print("output=" + output)
	if err != nil {
		log.Printf("[Specstest] Mount test filesystemtype = %v unsupported, err = %v...", fsname, err)
		resultString = "unsupported"
	} else if strings.EqualFold(strings.TrimSpace(output), "FilesystemSupportValidationContent") {
		resultString = "pass"
	} else {
		resultString = "failed"
	}
	return resultString, err
}

func main() {
	guestProgrammeFile := "mount_guest"
	outputFile := "spec_mount_out"
	err := hostsetup.SetupEnv(guestProgrammeFile, outputFile)
	if err != nil {
		log.Fatalf("[Specstest] Mount test filesystemtype: hostsetup.SetupEnv error, %v", err)
	}
	fmt.Println("Host enviroment setting up for runc is already!")

	testResult := new(TestResult)
	testResult.Mount.Type = make(map[string]string)

	// Filesytstem-tmpfs test
	testfsname := "tmpfs"
	if checkHostSupport(testfsname) {
		resultString, err := mountSupportTest(testfsname, testfsname, "/mountTest", "")
		if err != nil {
			log.Fatalf("[Specstest] Mount test filesystemtype= %v: mountSupportTest error, %v", testfsname, err)
		} else {
			testResult.Mount.Type[testfsname] = resultString
		}
	} else {
		testResult.Mount.Type[testfsname] = "Test host OS don't support the filesystem'"
	}

	//covert output struct to json format and output
	jsonString, err := json.Marshal(testResult)
	if err != nil {
		log.Fatalf("[Specstest] root test ouput json failed, err = %v...", err)
	}
	hostsetup.HostOutput(outputFile, string(jsonString))
}
