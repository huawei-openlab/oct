package factory

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

type RKT struct {
	name string
}

func (this *RKT) SetRT(runtime string) {
	this.name = "rkt"
}

func (this *RKT) GetRT() string {
	return "rkt"
}

func (this *RKT) NeedConvert() bool {
	return true
}

func (this *RKT) Convert(arg string, workingDir string) (string, error) {
	var cmd *exec.Cmd
	aciName := arg + ".aci"
	cmd = exec.Command("../plugins/oci2aci", "--debug", arg, aciName)
	cmd.Dir = workingDir //"./bundles"
	// cmd.stdin = os.Stdin
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("stderr err %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("stdout err %v", err)
	}

	var retStr string
	err = cmd.Start()
	if err != nil {
		retb, _ := ioutil.ReadAll(stderr)
		retStr = string(retb)
	} else {
		retb, _ := ioutil.ReadAll(stdout)
		retStr = string(retb)
	}

	return retStr, err
}

func (this *RKT) StartRT(specDir string) (string, error) {

	logrus.Debugf("Launcing runtime")
	/*rkt run 3.aci --interactive --insecure-skip-verify --mds-register=false --volume proc,kind=host,source=/bin --volume dev,kind=host,source=/bin --volume devpts,kind=host,source=/bin --volume shm,kind=host,source=/bin --volume mqueue,kind=host,source=/bin --volume sysfs,kind=host,source=/bin --volume cgroup,kind=host,source=/bin*/
	aciName := filepath.Base(specDir) + ".aci"
	aciPath := filepath.Dir(specDir)
	cmd := exec.Command("rkt", "run", aciName, "--interactive", "--insecure-skip-verify", "--mds-register=false",
		"--volume", "proc,kind=host,source=/bin", "--volume", "dev,kind=host,source=/bin", "--volume", "devpts,kind=host,source=/bin",
		"--volume", "shm,kind=host,source=/bin", "--volume", "mqueue,kind=host,source=/bin",
		"--volume", "sysfs,kind=host,source=/bin", "--volume", "cgroup,kind=host,source=/bin")
	cmd.Dir = aciPath
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	logrus.Debugf("Command done")

	if err != nil {
		return string(out), errors.New(string(out) + err.Error())
	}
	return string(out), nil
	/*if string(out) != "" {
		logrus.Infof("container output=%s\n", out)
	} else {
		logrus.Debugf("container output= nil\n")
	}
	if err != nil {
		return err
	}
	return nil*/
}

func (this *RKT) StopRT() error {
	return nil
}
