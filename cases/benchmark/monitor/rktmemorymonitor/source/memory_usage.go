package main

import (
	"../../source/adaptor"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/cadvisor/client"
	info "github.com/google/cadvisor/info/v1"
	"io/ioutil"
	"os"
	"strings"
)

type Memory struct {
	Container string `json:"container"`
	// Current memory usage, this includes all memory regardless of when it was
	// accessed.
	// Units: Bytes.
	Usage uint64 `json:"usage"`

	// The amount of working set memory, this includes recently accessed memory,
	// dirty memory, and kernel memory. Working set is <= "usage".
	// Units: Bytes.
	WorkingSet uint64 `json:"working_set"`

	ContainerData info.MemoryStatsMemoryData `json:"container_data,omitempty"`
	//HierarchicalData MemoryStatsMemoryData `json:"hierarchical_data,omitempty"`
}

func getContainerInfo(client *client.Client, container string) (containerInfo info.ContainerInfo, err error) {
	query := info.ContainerInfoRequest{}
	query.NumStats = 1
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

//container memory usage is current memory usage, this includes all memory regardless of when it was accessed.
// container memory working set  is the amount of working set memory, this includes recently accessed memory,
// dirty memory, and kernel memory. Working set is <= "usage".
func getMemory(cinfo info.ContainerInfo) (jsonString string, err error) {
	var temp string = ""
	for _, cstat := range cinfo.Stats {
		mem := Memory{Container: cinfo.Name,
			Usage:         cstat.Memory.Usage,
			WorkingSet:    cstat.Memory.WorkingSet,
			ContainerData: cstat.Memory.ContainerData,
		}

		tempJsonString, err := json.Marshal(mem)
		if err != nil {
			return "", err
		}

		temp = temp +string(tempJsonString)
	}
	return temp, nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("commad must has one parameters!\n")
		return
	}
	var testingProject = os.Args[1] //"docker"  or  "rkt"
	if testingProject != "docker" && testingProject != "rkt" {
		fmt.Printf("commad is %v %v, is not corrected!\n", os.Args[0], os.Args[1])
		return
	}

	var containers []string
	client, err := client.NewClient("http://localhost:8080/")
	if err != nil {
		fmt.Printf("tried to make client and got error %v\n", err)
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
		fmt.Printf("getContainerName fail, error: %v\n", err)
		return
	}

	var jsonString string = ""
	for _, container := range containers {
		//fmt.Printf("container %v's memory info: \n", container)
		cinfo, err := getContainerInfo(client, container)
		if err != nil {
			fmt.Printf("getContainerInfo fail and got error %v\n", err)
			return
		}
		temp, err := getMemory(cinfo)
		if err != nil {
			fmt.Printf("getMemory faile, error: %v\n", err)
		}
		jsonString = jsonString + temp
	}

	err = ioutil.WriteFile("./"+testingProject+"_memory.json",  []byte(jsonString), 0666)
	if err != nil {
		fmt.Printf("ioutil.WriteFile faile, error: %v\n", err)
	}
}
