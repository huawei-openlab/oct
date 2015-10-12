package utils

import (
	"log"
	"os"
)

func SetRight(file string, uid uint32, gid uint32) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Open file %v error %v", file, err)
	}
	defer f.Close()

	err = f.Chown(int(uid), int(gid))
	if err != nil {
		log.Fatalf("Chown file %v error %v", file, err)
	}

	// Read all files under file
	ff, _ := f.Readdirnames(0)
	for _, fi := range ff {
		subFile := file + "/" + fi
		err = os.Chown(subFile, int(uid), int(gid))
		if err != nil {
			log.Fatalf("Chown file %v error %v", subFile, err)
		}
	}
}
