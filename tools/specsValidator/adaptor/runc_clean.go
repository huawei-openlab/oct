package adaptor

import (
	"log"
	"os/exec"
)

func CleanRunc() {
	cmd := exec.Command("/bin/bash", "-c", "kill -9 `ps -ef|grep runc|grep -v grep|awk '{print $2}'`")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[clean runc] kill process error , %v", err)
	}
	cmd = exec.Command("/bin/bash", "-c", "rm -r /run/oci/specsValidator")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("[clean runc] delete folder error , %v", err)
	}
	log.Println("clean runc success")

}
