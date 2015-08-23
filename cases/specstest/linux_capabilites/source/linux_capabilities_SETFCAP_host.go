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
	specs "./../../source/specs"
	"fmt"
	"log"
)

func linuxCapabilitiesTestSETFCAP() {

	//set file path
	configjsonFilePath := "./../../source/config.json"
	guestProgrammeFileName := "linux_capabilities_SETFCAP_guest"
	outputFileName := "linux_capabilities_SETFCAP"

	//setup the guest enviroment
	err := hostsetup.SetupEnv(guestProgrammeFileName, outputFileName)
	if err != nil {
		log.Fatalf("[Specstest] linux Capabilities SETFCAP test: hostsetup.SetupEnv error, %v", err)
	}
	fmt.Println("Host enviroment setting up for runc is already!")

	//read the config.json and edit and convert
	var linuxspec *specs.LinuxSpec
	linuxspec, err = configconvert.ConfigToLinuxSpec(configjsonFilePath)
	if err != nil {
		log.Fatalf("[Specstest] linux Capabilities SETFCAP test: reading config error, %v", err)
	}
	linuxspec.Spec.Root.Path = "./rootfs_rootconfig"
	linuxspec.Process.Args = []string{("./" + guestProgrammeFileName)}
	capbilitestest := []string{"SETFCAP"}
	linuxspec.Linux.Capabilities = capbilitestest
	err = configconvert.LinuxSpecToConfig(configjsonFilePath, linuxspec)
	if err != nil {
		log.Fatalf("[Specstest] linux Capabilities SETFCAP test:writing config error, %v", err)
	}
	fmt.Println("Host enviroment for runc is already!")

}

func main() {
	linuxCapabilitiesTestSETFCAP()
}
