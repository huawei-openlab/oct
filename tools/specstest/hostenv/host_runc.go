package hostenv

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//https://github.com/opencontainers/specs.git
func cloneDeps(path string, repo string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("CloneDeps create path err = %v", err)
	}

	cloneString := "https://github.com/opencontainers/" + repo + ".git"
	cmd := exec.Command("git", "clone", cloneString)
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = path
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("CloneDeps git clone %v err", cloneString)
	}
	return nil
}

func upateRev(Rev string, repo string) (string, error) {
	fmt.Printf("Rev : %v", Rev)
	goPath := os.Getenv("GOPATH")
	path := goPath + "/src/github.com/opencontainers/" + repo
	pPath := goPath + "/src/github.com/opencontainers/"
	if !exists(path) {
		err := cloneDeps(pPath, repo)
		if err != nil {
			return path, err
		}
	}

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
	err = setRuncEnv()
	return err
}

func setRuncEnv() error {
	runcPath := "/usr/local/bin/"
	result := os.Getenv("PATH")
	if !strings.Contains(result, runcPath) {
		newPath := runcPath + ":" + result
		if err := os.Setenv("PATH", newPath); err != nil {
			log.Fatalf("SetRuncEnv err %v", err)
		}
	}
	return nil
}
