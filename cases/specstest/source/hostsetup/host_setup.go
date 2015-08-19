package hostsetup

import (
	"fmt"
	"log"
	"os/exec"
)

func SetupEnv(guestFile string, outputFile string) error {
	var cmd *exec.Cmd
	var err error

	cmd = exec.Command("/bin/sh", "-c", "mkdir -p /tmp/testtool")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: mkdir testtool dir error, %v", err)
	}

	if guestFile != "" {
		fmt.Println("Build guest programme...................")
		cmd = exec.Command("/bin/sh", "-c", "go build "+guestFile+".go")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: build guest programme error, %v", err)
		}
		cmd = exec.Command("/bin/sh", "-c", "mv "+guestFile+" /tmp/testtool/")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: mv guest programme error, %v", err)
		}
	}

	if outputFile != "" {
		fmt.Println("Touch output file...................")
		cmd = exec.Command("/bin/sh", "-c", "touch /tmp/testtool/"+outputFile+".json")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: touch %s.json err, %v", outputFile, err)
		}

		cmd = exec.Command("/bin/sh", "-c", "chown root:root  /tmp/testtool/"+outputFile+".json")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: change to root power err, %v", err)
		}

	}
	fmt.Println("Pull docker image...................")
	cmd = exec.Command("/bin/sh", "-c", "docker pull ubuntu:14.04")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: pull image error, %v", err)
	}
	fmt.Println("Export docker filesystem...................")
	cmd = exec.Command("/bin/sh", "-c", "docker export $(docker create ubuntu) > ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: export image error, %v", err)
	}
	cmd = exec.Command("/bin/sh", "-c", "mkdir -p ./../../source/rootfs_rootconfig")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: create rootfs dir error, %v", err)
	}
	fmt.Println("Extract rootfs...................")
	cmd = exec.Command("/bin/sh", "-c", "tar -C ./../../source/rootfs_rootconfig -xf ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: create rootfs content error, %v", err)
	}

	return nil
}
