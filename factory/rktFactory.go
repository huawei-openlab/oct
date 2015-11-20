package factory

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/huawei-openlab/oct/hooks"
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

func (this *RKT) PreStart(configArgs string) error {
	return nil
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

func (this *RKT) PostStart(configArgs string, containerout string) error {
	if strings.Contains(configArgs, "-args=./runtimetest --args=vna") {
		if err := hooks.NamespacePostStart(containerout); err != nil {
			return nil
		}
	}
	return nil
}

func (this *RKT) StopRT() error {
	return nil
}
