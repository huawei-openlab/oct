package hostenv

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func upateRev(Rev string, repo string) (string, error) {
	fmt.Println("Rev : %v", Rev)
	goPath := os.Getenv("GOPATH")
	path := goPath + "/src/github.com/opencontainers/" + repo
	if !exists(path) {
		log.Fatalf("In hostenv path: %v is not exist", path)
	}

	/*	err := os.Chdir(path)
		if err != nil {
			log.Fatalf("In Hostenv chdir to path : %v err", path)
		}*/

	cmd := exec.Command("git", "checkout", Rev)
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = path
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("UpateRev git checkout %v err", Rev)
	}

	return path, nil
}

func UpateRuncRev(runcRev string) error {
	path, _ := upateRev(runcRev, "runc")

	cmd := exec.Command("make")
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = path
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("UpateRuncRev make runc err : %v", err)
	}

	cmd = exec.Command("make", "install")
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = path
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("UpateRuncRev make install runc err : %v", err)
	}
	return err
}
