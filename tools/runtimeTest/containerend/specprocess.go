package main

import (
	// "fmt"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	// "io"
	"fmt"
	// "strings"
)

func main() {
	procFile := "/proc/self/status"
	suid := utils.GetJob("Uid", procFile)
	guid := utils.GetJob("Gid", procFile)
	groups := utils.GetJob("Groups", procFile)
	suid = suid + guid + groups
	fmt.Println(suid)
}
