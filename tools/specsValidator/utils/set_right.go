package utils

import (
	"log"
	"os"
)

func SetRight(file string, uid int32, gid int32) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Open file %v error %v", file, err)
	}
	defer f.Close()
	err = f.Chown(int(uid), int(gid))
	if err != nil {
		log.Fatalf("Chown file %v error %v", file, err)
	}
}
