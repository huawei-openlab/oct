package factory

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/huawei-openlab/oct/hooks"
)

type Runc struct {
	name string
}

func (this *Runc) SetRT(runtime string) {
	this.name = "runc"
}

func (this *Runc) GetRT() string {
	return "runc"
}

func (this *Runc) PreStart(configArgs string) error {
	return nil
}

func (this *Runc) StartRT(specDir string) (string, error) {
	logrus.Debugf("Launcing runtime")

	cmd := exec.Command("runc", "start")
	cmd.Dir = specDir
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	/*if err := hostendvalidate.ContainerOutputHandler(string(out)); err != nil {
		return err
	}*/
	logrus.Debugf("Command done")
	if err != nil {
		return string(out), errors.New(string(out) + err.Error())
	}

	return string(out), nil

	/*if string(out) != "" {

		logrus.Printf("container output=%s\n", out)
	} else {
		logrus.Debugf("container output= nil\n")
	}
	if err != nil {
		return err
	}
	return nil*/
}

func (this *Runc) PostStart(configArgs string, containerout string) error {
	if strings.Contains(configArgs, "-args=./runtimetest --args=vna") {
		if err := hooks.NamespacePostStart(containerout); err != nil {
			return nil
		}
	}
	return nil
}

func (this *Runc) StopRT() error {
	return nil
}
