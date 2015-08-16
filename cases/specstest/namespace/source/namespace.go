// Copyright 2014 Google Inc. All Rights Reserved.
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

package main

import (
	"../../source/configconvert"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type TestResult struct {
	Pid     string `json:"Linuxspec.Spec.Namespaces.pid"`
	Mount   string `json:"Linuxspec.Spec.Namespaces.mount"`
	Ipc     string `json:"Linuxspec.Spec.Namespaces.ipc"`
	Network string `json:"Linuxspec.Spec.Namespaces.netwok"`
	Uts     string `json:"Linuxspec.Spec.Namespaces.uts"`
}

func pullImage() {
	fmt.Println("Pull docker image...................")
	cmd := exec.Command("/bin/sh", "-c", "docker pull ubuntu:14.04")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: pull image error, %v", err)
	}
	fmt.Println("Export docker filesystem...................")
	cmd = exec.Command("/bin/sh", "-c", "docker export $(docker create ubuntu) > ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: export image error, %v", err)
	}
	cmd = exec.Command("/bin/sh", "-c", "mkdir -p ./rootfs_rootconfig")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: create rootfs_rootconfig dir error, %v", err)
	}
	fmt.Println("Extract rootfs...................")
	cmd = exec.Command("/bin/sh", "-c", "tar -C ./rootfs_rootconfig -xf ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: create rootfs_rootconfig content error, %v", err)
	}

	//for adptor the ../../source/config.json
	cmd = exec.Command("/bin/sh", "-c", "mkdir -p /tmp/testtool")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: mkdir testtool dir error, %v", err)
	}
}

func getContainerNS(configFilePath string) ([]string, error) {
	//1. Run runc to start container, and run 1.sh in container to get namespace in container
	cmd := exec.Command("runc", configFilePath)
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = "./"
	/*
		return namespace's format like
		lrwxrwxrwx 1 root root 0  8月 13 09:36 ipc -> ipc:[4026531839]
		lrwxrwxrwx 1 root root 0  8月 13 09:38 mnt -> mnt:[4026531840]
		lrwxrwxrwx 1 root root 0  8月 13 09:38 net -> net:[4026531956]
		lrwxrwxrwx 1 root root 0  8月 13 09:38 pid -> pid:[4026531836]
		lrwxrwxrwx 1 root root 0  8月 13 09:38 user -> user:[4026531837]
		lrwxrwxrwx 1 root root 0  8月 13 09:38 uts -> uts:[4026531838]
	*/
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.New(string(out) + err.Error())
	}

	str := strings.TrimSuffix(string(out), "\n")
	str = strings.TrimSpace(str)
	if str == "" {
		return nil, errors.New("can not find namespace in container.")
	}

	lines := strings.Split(str, "\n")
	var cNamespaces []string
	for i, line := range lines {
		//remove first line: total 0\n
		if 0 == i {
			continue
		}
		tmp := strings.Split(line, "->")
		cNamespaces = append(cNamespaces, strings.TrimSpace(tmp[1]))
	}
	return cNamespaces, nil
}

func getHostNS(configNSPath bool) (string, error) {
	//2. Get the namespaces on host
	/*
		ipc:[4026531839]
		mnt:[4026531840]
		mnt:[4026531856]
		mnt:[4026532129]
		net:[4026531956]
		net:[4026532265]
		pid:[4026531836]
		pid:[4026532263]
		user:[4026531837]
		user:[4026532227]
		uts:[4026531838]
	*/
	var cmd *exec.Cmd
	if configNSPath {
		cmd = exec.Command("/bin/sh", "-c", "readlink /proc/1/ns/* | sort -u")
	} else {
		cmd = exec.Command("/bin/sh", "-c", "readlink /proc/*/ns/* | sort -u")
	}

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	str := strings.TrimSuffix(string(out), "\n")
	str = strings.TrimSpace(str)
	if str == "" {
		return "", errors.New("can not find namespace on host.")
	}
	return str, nil
}

func testNSConfig() {

	//1. set the config.json
	fmt.Println("Host enviroment setting up for runc is already!")

	inFilePath := "../../source/config.json"
	outFilePath := "./configns.json"

	linuxspec, err := configconvert.ConfigToLinuxSpec(inFilePath)
	if err != nil {
		log.Fatalf("Specstestroot readonly test: readconfig error, %v", err)
	}

	//namespace configed with nil
	for i, _ := range linuxspec.Linux.Namespaces {
		switch linuxspec.Linux.Namespaces[i].Type {
		case "pid":
			linuxspec.Linux.Namespaces[i].Path = "/proc/1/ns/pid"
		case "mount":
			linuxspec.Linux.Namespaces[i].Path = ""
		case "ipc":
			linuxspec.Linux.Namespaces[i].Path = "/proc/1/ns/ipc"
		case "uts":
			linuxspec.Linux.Namespaces[i].Path = "/proc/1/ns/uts"
		case "network":
			linuxspec.Linux.Namespaces[i].Path = "/proc/1/ns/net"
		}
	}

	linuxspec.Spec.Process.Args[0] = "ls"
	linuxspec.Spec.Process.Args = append(linuxspec.Spec.Process.Args, "-l", "/proc/self/ns")

	err = configconvert.LinuxSpecToConfig(outFilePath, linuxspec)
	if err != nil {
		log.Fatalf("Specstest root readonly test: writeconfig error, %v", err)
	}
	fmt.Println("Host enviroment for runc is already!")

	//2. Run runc to start container, and  get namespace in container
	cNamespaces, err := getContainerNS(outFilePath)
	if err != nil {
		log.Fatalf("getContainerNS error, %v", err)
	}
	//3. Get the namespaces on host
	str, err := getHostNS(true)
	if err != nil {
		log.Fatalf("getHostNS error, %v", err)
	}

	//4. judge and output the result
	for i, _ := range linuxspec.Linux.Namespaces {
		//config.json's "mount" and "network" is  different with os's "mnt" and "net"
		if linuxspec.Linux.Namespaces[i].Type == "mount" {
			linuxspec.Linux.Namespaces[i].Type = "mnt"
		}
		if linuxspec.Linux.Namespaces[i].Type == "network" {
			linuxspec.Linux.Namespaces[i].Type = "net"
		}
	}

	var result TestResult
	for _, ns := range cNamespaces {
		for _, nstmp := range linuxspec.Linux.Namespaces {
			if strings.Contains(ns, nstmp.Type) {
				ts := strings.Split(ns, ":")
				if strings.Contains(str, ns) {
					switch ts[0] {
					case "pid":
						result.Pid = "pass"
					case "ipc":
						result.Ipc = "pass"
					case "mnt":
						result.Mount = "pass"
					case "net":
						result.Network = "pass"
					case "uts":
						result.Uts = "pass"
					}
				} else {
					switch ts[0] {
					case "pid":
						result.Pid = "failed"
					case "ipc":
						result.Pid = "failed"
					case "mnt":
						result.Mount = "failed"
					case "net":
						result.Network = "failed"
					case "uts":
						result.Uts = "failed"
					}
				}
			}
		}
	}

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Fatalf(" json.Marshal error, %v", err)
	}

	outFile := "./namespace_config_out.json"
	ioutil.WriteFile(outFile, []byte(jsonString), 666)
}

func testNSNotConfig() {
	fmt.Println("Host enviroment setting up for runc is already!")

	inFilePath := "../../source/config.json"
	outFilePath := "./confignotns.json"

	linuxspec, err := configconvert.ConfigToLinuxSpec(inFilePath)
	if err != nil {
		log.Fatalf("Specstestroot readonly test: readconfig error, %v", err)
	}

	//namespace configed with nil
	for _, ns := range linuxspec.Linux.Namespaces {
		ns.Path = ""
	}

	linuxspec.Spec.Process.Args[0] = "ls"
	linuxspec.Spec.Process.Args = append(linuxspec.Spec.Process.Args, "-l", "/proc/1/ns")

	err = configconvert.LinuxSpecToConfig(outFilePath, linuxspec)
	if err != nil {
		log.Fatalf("Specstest root readonly test: writeconfig error, %v", err)
	}
	fmt.Println("Host enviroment for runc is already!")

	//2. Run runc to start container, and  get namespace in container
	cNamespaces, err := getContainerNS(outFilePath)
	if err != nil {
		log.Fatalf("getContainerNS error, %v", err)
	}
	//3. Get the namespaces on host
	str, err := getHostNS(false)
	if err != nil {
		log.Fatalf("getHostNS error, %v", err)
	}

	//4. judge and output the result
	for i, _ := range linuxspec.Linux.Namespaces {
		//config.json's "mount" and "network" is  different with os's "mnt" and "net"
		if linuxspec.Linux.Namespaces[i].Type == "mount" {
			linuxspec.Linux.Namespaces[i].Type = "mnt"
		}
		if linuxspec.Linux.Namespaces[i].Type == "network" {
			linuxspec.Linux.Namespaces[i].Type = "net"
		}
	}

	var result TestResult
	for _, ns := range cNamespaces {
		for _, nstmp := range linuxspec.Linux.Namespaces {
			if strings.Contains(ns, nstmp.Type) {
				ts := strings.Split(ns, ":")
				if strings.Contains(str, ns) {
					switch ts[0] {
					case "pid":
						result.Pid = "failed"
					case "ipc":
						result.Ipc = "failed"
					case "mnt":
						result.Mount = "failed"
					case "net":
						result.Network = "failed"
					case "uts":
						result.Uts = "failed"
					}
				} else {
					switch ts[0] {
					case "pid":
						result.Pid = "pass"
					case "ipc":
						result.Pid = "pass"
					case "mnt":
						result.Mount = "pass"
					case "net":
						result.Network = "pass"
					case "uts":
						result.Uts = "pass"
					}
				}
			}
		}
	}

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Fatalf(" json.Marshal error, %v", err)
	}

	outFile := "./namespace_notconfig_out.json"
	ioutil.WriteFile(outFile, []byte(jsonString), 666)
}

func main() {
	pullImage()
	testNSNotConfig()
	testNSConfig()
}
