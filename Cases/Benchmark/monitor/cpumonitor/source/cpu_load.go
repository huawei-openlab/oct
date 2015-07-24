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
	"fmt"
	"github.com/google/cadvisor/client"
	info "github.com/google/cadvisor/info/v1"
	"log"
	"os/exec"
)

func cpuTotalUsage() {
	staticClient, err := client.NewClient("http://localhost:8080/")
	if err != nil {
		log.Fatalf("CpuUsage monitor tried to make client and got error %v", err)
		return
	}

	//exec shell in host machine to get docker container id
	cmd := exec.Command("/bin/sh", "-c", "docker ps  -q")
	short_id, err := cmd.Output()

	cmd = exec.Command("/bin/sh", "-c", "docker inspect -f   '{{.Id}}' "+string(short_id))
	full_id, err := cmd.Output()
	fmt.Println("Cpu monitor find out the docker container ID %v", string(full_id))

	//containerName := "/docker/container id"
	containerName := "/docker/" + string(full_id)
	query := &info.ContainerInfoRequest{}

	//get ContainerInfo structure according the client
	cInfo, err := staticClient.ContainerInfo(containerName, query)
	if err != nil {
		log.Fatalf("Cpu monitor get ContainerInfo err %v", err)
		return
	}
	fmt.Println("Cpu monitor get container Info container name: %v ,now start get cpu load ", cInfo.Name)

	for i := 1; i < len(cInfo.Stats); i++ {
		cur := cInfo.Stats[i]
		fmt.Println("Container cpu load : ", float64(cur.Cpu.LoadAverage)/1000)
	}
}

func main() {
	cpuTotalUsage()
}
