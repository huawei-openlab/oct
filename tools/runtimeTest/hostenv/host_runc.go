package hostenv

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"log"
	"os"
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

	_, err = utils.ExecCmdStr("git", path, "clone", cloneString)
	/*cmd := exec.Command("git", "clone", cloneString)
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = path
	_, err = cmd.Output()*/
	if err != nil {
		log.Fatalf("CloneDeps git clone %v err", cloneString)
	}
	return nil
}

func upateRev(Rev string, repo string) (string, error) {

	goPath := os.Getenv("GOPATH")
	path := goPath + "/src/github.com/opencontainers/" + repo
	pPath := goPath + "/src/github.com/opencontainers/"
	err := checkout(pPath, path, Rev, repo)
	if err != nil {
		return path, err
	}
	/*if !exists(path) {
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
	}*/

	// If it is specs, should also have to update in Godeps
	if repo == "specs" {
		path, err = os.Getwd()
		if err != nil {
			return path, err
		}
		// fmt.Println(path)
		pPath = getParentDirectory(path)
		// fmt.Println(pPath)
		path = pPath + "/Godeps/_workspace/src/github.com/opencontainers/specs/"
		pPath = pPath + "/Godeps/_workspace/src/github.com/opencontainers/"
		err := checkout(pPath, path, Rev, repo)
		if err != nil {
			return path, err
		}
	}

	return path, nil
}

func checkout(pPath string, path string, Rev string, repo string) error {
	if !exists(path) {
		err := cloneDeps(pPath, repo)
		if err != nil {
			return err
		}
	}

	_, err := utils.ExecCmdStr("git", path, "checkout", Rev)
	/*
		cmd := exec.Command("git", "checkout", Rev)
		cmd.Stderr = os.Stderr
		//cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Dir = path
		_, err := cmd.Output()*/
	if err != nil {
		return err
	}
	return nil
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func UpateRuncRev(runcRev string, tags string) error {
	path, _ := upateRev(runcRev, "runc")

	_, err := utils.ExecCmdStr("make", path, "BUILDTAGS="+tags)
	/*cmd := exec.Command("make", "BUILDTAGS="+tags)
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = path
	_, err := cmd.Output()*/
	if err != nil {
		log.Fatalf("UpateRuncRev make runc err : %v", err)
	}

	_, err = utils.ExecCmdStr("make", path, "install")
	/*cmd = exec.Command("make", "install")
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = path
	_, err = cmd.Output()*/
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
