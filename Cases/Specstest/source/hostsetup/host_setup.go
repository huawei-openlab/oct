package hostsetup

import (
	"log"
	"os/exec"
)

func SetupEnv(guestFile string) error {
	var cmd *exec.Cmd
	var err error
	if guestFile != "" {
		cmd = exec.Command("/bin/sh", "-c", "go build "+guestFile+".go")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: build guest programme error, %v", err)
		}

		cmd = exec.Command("/bin/sh", "-c", "mkdir -p /tmp/testtool")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: mkdir testtool dir error, %v", err)
		}
		cmd = exec.Command("/bin/sh", "-c", "mv "+guestFile+" /tmp/testtool/")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: mv guest programme error, %v", err)
		}

		cmd = exec.Command("/bin/sh", "-c", "touch  /tmp/testtool/readonly_true_out.txt")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: touch readonly_true_out.txt err, %v", err)
		}

		cmd = exec.Command("/bin/sh", "-c", "chown root:root  /tmp/testtool/readonly_true_out.txt")
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("Specstest root readonly test: change to root power err, %v", err)
		}
	}

	cmd = exec.Command("/bin/sh", "-c", "docker pull ubuntu:14.04")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: pull image error, %v", err)
	}
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

	cmd = exec.Command("/bin/sh", "-c", "tar -C ./../../source/rootfs_rootconfig -xf ubuntu.tar")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("Specstest root readonly test: create rootfs content error, %v", err)
	}

	return nil
}
