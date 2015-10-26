package utils

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func GetJob(job string, file string) string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Open file %v error %v", file, err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.Contains(line, job) {
			return line
		} else {
			continue
		}
	}
	return ""
}
