package main

import (
	// "fmt"
	"github.com/huawei-openlab/oct/tools/specsValidator/utils"
	// "io"
	"log"
	"os"
	// "strings"
)

func main() {

	fileName := "./containerend_out.txt"
	fout, err := os.Create(fileName)
	defer fout.Close()
	if err != nil {
		log.Fatalf("Create %v error %v", fileName, err)
	}

	procFile := "/proc/self/status"
	suid := utils.GetJob("Uid", procFile)
	guid := utils.GetJob("Gid", procFile)
	groups := utils.GetJob("Groups", procFile)
	suid = suid + guid + groups
	fout.WriteString(suid)
}
