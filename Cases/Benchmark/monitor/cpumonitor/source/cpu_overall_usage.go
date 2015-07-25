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

func cpuOverallUsage() {
	staticClient, err := client.NewClient("http://localhost:8080/")
	if err != nil {
		log.Fatalf("Cpu monitor tried to make client and got error %v", err)
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

	mInfo, err := staticClient.MachineInfo()
	if err != nil {
		log.Fatalf("Cpu monitor try to get MachineInfo and got err %v", err)
		return
	}

	fmt.Println("Cpu monitor get container Info container name: %v ,now start get overall cpu usage", cInfo.Name)

	cur := cInfo.Stats[len(cInfo.Stats)-1]

	if len(cInfo.Stats) >= 2 {
		prev := cInfo.Stats[len(cInfo.Stats)-2]
		rawUsage := float64(cur.Cpu.Usage.Total - prev.Cpu.Usage.Total)
		intervalInNs := float64((cur.Timestamp).Sub(prev.Timestamp).Nanoseconds())
		cpuUsage := ((rawUsage / intervalInNs) / float64(mInfo.NumCores)) * 100
		fmt.Printf("cpuUsage %.02f \n", cpuUsage)
	}

}

func main() {
	cpuOverallUsage()
}
