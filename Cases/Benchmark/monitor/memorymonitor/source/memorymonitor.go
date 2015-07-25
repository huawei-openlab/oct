package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/google/cadvisor/client"
	info "github.com/google/cadvisor/info/v1"
)

func GetAllContainer(client *client.Client) ([]info.ContainerInfo, error) {
	query := info.ContainerInfoRequest{}
	query.NumStats = 1
	cinfos, err := client.SubcontainersInfo("/", &query)
	return cinfos, err
}

//container memory usage is current memory usage, this includes all memory regardless of when it was accessed.
// Units: Bytes.

func GetContainerMemoryUsage(cinfos []info.ContainerInfo) {

	fmt.Printf("container memeory usage:\n")
	for _, cinfo := range cinfos {
		for _, cstat := range cinfo.Stats {
			fmt.Printf("container name: %s,  usage:%d\n", cinfo.Name, cstat.Memory.Usage)
		}
	}
}

// container memory working set  is the amount of working set memory, this includes recently accessed memory,
// dirty memory, and kernel memory. Working set is <= "usage".
// Units: Bytes.

func GetContainerMemoryWorkingSet(cinfos []info.ContainerInfo) {
	fmt.Printf("container memeory  working set:\n")
	for _, cinfo := range cinfos {
		for _, cstat := range cinfo.Stats {
			fmt.Printf("container name: %s,  workset bytes:%d\n", cinfo.Name, cstat.Memory.WorkingSet)
		}
	}
}

func main() {
	client, err := client.NewClient("http://localhost:8080/")
	if err != nil {
		glog.Errorf("tried to make client and got error %v", err)
		return
	}

	cinfos, err := GetAllContainer(client)
	if err != nil {
		glog.Errorf("tried to SubcontainersInfo and got error %v", err)
		return
	}
	GetContainerMemoryUsage(cinfos)
	fmt.Printf("\n")
	GetContainerMemoryWorkingSet(cinfos)

}
