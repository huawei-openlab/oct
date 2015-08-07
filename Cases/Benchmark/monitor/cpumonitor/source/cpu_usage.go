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
	adaptor "./../../source/adaptor"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/cadvisor/client"
	info "github.com/google/cadvisor/info/v1"
	"io/ioutil"
	"log"
	"os"
	"strings"
	//"os/exec"
)

// CPU usage time statistics.
type CpuUsageInfo struct {
	Usage       CpuUsage `json:"cpu_usage"`
	ContainerID string   `json:"container_id"`
}

type CpuUsage struct {
	// Total CPU usage.
	// Units: nanoseconds
	TotalUsage float64 `json:"total_usage"`

	// Per CPU/core usage of the container.
	// Unit: nanoseconds.
	PerCoreUsage []CpuPerCoreUsage `json:"percore_usage"`

	Load float64 `json:"load"`

	OverallUsage float64 `json:"overall_usage"`

	BreakdownUsage CpuBreakdown `json:"breakdown_usage"`
}

type CpuBreakdown struct {
	// Time spent in user space.
	// Unit: nanoseconds
	UserUsage float64 `json:"user_usage"`

	// Time spent in kernel space.
	// Unit: nanoseconds
	SystemUsage float64 `json:"system_usage"`
}

type CpuPerCoreUsage struct {
	coreID       int     `json:"core_ID"`
	percoreUsage float64 `json:"percore_usage"`
}

func getCpu(cInfo info.ContainerInfo, mInfo *info.MachineInfo, cpuusageinfo *CpuUsageInfo) (err error) {
	cur := cInfo.Stats[len(cInfo.Stats)-1]
	fmt.Println(cpuusageinfo.ContainerID)
	if len(cInfo.Stats) >= 2 {
		prev := cInfo.Stats[len(cInfo.Stats)-2]
		rawUsage := float64(cur.Cpu.Usage.Total - prev.Cpu.Usage.Total)
		intervalInNs := float64((cur.Timestamp).Sub(prev.Timestamp).Nanoseconds())
		cpuusageinfo.Usage.OverallUsage = ((rawUsage / intervalInNs) / float64(mInfo.NumCores)) * 100
		fmt.Printf("Container cpu overall usage %.02f \n", cpuusageinfo.Usage.OverallUsage)
	}
	for i := 1; i < len(cInfo.Stats); i++ {
		cur := cInfo.Stats[i]
		prev := cInfo.Stats[i-1]
		//get interval time duration between the two sample
		f := float64((cur.Timestamp).Sub(prev.Timestamp).Nanoseconds())
		cpuusageinfo.Usage.Load = float64(cur.Cpu.LoadAverage) / 1000
		cpuusageinfo.Usage.TotalUsage = float64(cur.Cpu.Usage.Total-prev.Cpu.Usage.Total) / f
		cpuusageinfo.Usage.BreakdownUsage.SystemUsage = float64(cur.Cpu.Usage.User-prev.Cpu.Usage.User) / f
		cpuusageinfo.Usage.BreakdownUsage.UserUsage = float64(cur.Cpu.Usage.System-prev.Cpu.Usage.System) / f

	}
	return nil
}

func getContainerInfo(client *client.Client, container string) (containerInfo info.ContainerInfo, err error) {
	query := info.ContainerInfoRequest{}
	cinfos, err := client.SubcontainersInfo("/", &query)
	if err != nil {
		return info.ContainerInfo{}, err
	}
	tempContainer := "/" + container
	for _, cinfo := range cinfos {
		if strings.HasSuffix(cinfo.Name, tempContainer) {
			return cinfo, nil
		}
	}
	return info.ContainerInfo{}, errors.New("not find container " + container)
}
func main() {

	if len(os.Args) < 2 {
		log.Fatalf("commad must has one parameters!\n")
		return
	}
	var testingProject = os.Args[1] //"docker"  or  "rkt"
	if testingProject != "docker" && testingProject != "rkt" {
		log.Fatalf("commad is %v %v, is not corrected!\n", os.Args[0], os.Args[1])
		return
	}

	var containers []string
	client, err := client.NewClient("http://localhost:8080/")
	if err != nil {
		log.Fatalf("tried to make client and got error %v\n", err)
		return
	}

	switch testingProject {
	case "docker":
		containers, err = adaptor.GetDockerContainers()
	case "rkt":
		containers, err = adaptor.GetRktContainers()
	default:
		return
	}
	if err != nil {
		log.Fatalf("getContainerName fail, error: %v\n", err)
		return
	}
	//cpuusageinfo := new(CpuUsageInfo)
	mInfo, err := client.MachineInfo()
	var jsonString []byte
	for _, container := range containers {
		fmt.Printf("container %v's cpu info: \n", container)
		cInfo, err := getContainerInfo(client, container)
		if err != nil {
			fmt.Printf("getContainerInfo fail and got error %v\n", err)
			return
		}
		cpuUsageInfo := new(CpuUsageInfo)
		cpuUsageInfo.ContainerID = cInfo.Name
		//fmt.Println(cpuUsageInfo.ContainerID)
		//var usageInfo CpuUsageInfo
		err = getCpu(cInfo, mInfo, cpuUsageInfo)
		fmt.Println(cpuUsageInfo)

		//fmt.Println(usageInfo.ContainerID)
		jsonString, err = json.Marshal(cpuUsageInfo)
		if err != nil {
			log.Fatalf("convert to json err, error:  %v\n", err)
			return
		}

	}
	err = ioutil.WriteFile("./"+testingProject+"_cpu.json", []byte(jsonString), 0666)

}
