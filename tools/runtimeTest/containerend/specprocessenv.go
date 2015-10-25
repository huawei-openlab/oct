package main

import (
	"fmt"
	"os"
	// "strings"
	"io/ioutil"
	"log"
)

func main() {
	procFile := "/proc/self/environ"
	fd, err := os.Open(procFile)
	if err != nil {
		log.Fatalf("Open procFile error: %v", err)
	}

	s, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Fatalf("Read procFile error: %v", err)
	}
	fmt.Println(string(s))
	//fmt.Println("specprocessenv run sucessful with env set")
}
