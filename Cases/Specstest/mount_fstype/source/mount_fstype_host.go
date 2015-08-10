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

package main

import (
	configconvert "./../../source/configconvert"
	hostsetup "./../../source/hostsetup"
	"fmt"
	specs "github.com/opencontainers/specs"
	"log"
)

func setFsListAndExcuteTest() {
	// filesystemlists := []string{"ext2", "ext3", "ext4", "ntfs", "tmpfs", "sysfs"}
	filesystemlists := []string{"ext2"}
	for _, fs := range filesystemlists {
		testFsTypeSupport(fs)
	}
}

func testFsTypeSupport(fstest string) {
	//set file path
	configjsonFilePath := "./../../source/config.json"
	guestProgrammeFileName := ""
	outputFileName := "mount_fstypesupport_out"
	//setup the guest enviroment
	err := hostsetup.SetupEnv(guestProgrammeFileName, outputFileName)
	if err != nil {
		log.Fatalf("[Specstest] mount filesystem support test: hostsetup.SetupEnv error, %v", err)
	}
	fmt.Println("Host enviroment setting up for runc is already!")

	//read the config.json and edit and convert
	var linuxspec *specs.LinuxSpec
	linuxspec, err = configconvert.ConfigToLinuxSpec(configjsonFilePath)
	if err != nil {
		log.Fatalf("[Specstest] mount filesystem support test: reading config error, %v", err)
	}
	linuxspec.Spec.Root.Path = "./rootfs_rootconfig"
	mountsorigin := specs.Mount{"proc", "proc", "/proc", ""}
	mountsadd := specs.Mount{fstest, "/tmp/test", "/testfs", ""}
	mountsnew := []specs.Mount{mountsorigin, mountsadd}
	linuxspec.Mounts = mountsnew
	err = configconvert.LinuxSpecToConfig(configjsonFilePath, linuxspec)
	if err != nil {
		log.Fatalf("[Specstest] mount filesystem support test:writing config error, %v", err)
	}
	fmt.Println("Host enviroment for runc is already!")

}

func main() {
	setFsListAndExcuteTest()
}
