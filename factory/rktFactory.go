package factory

import (
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
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

func (this *RKT) StartRT(specDir string) error {

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
	if string(out) != "" {
		logrus.Debugf("container output=%s\n", out)
	} else {
		logrus.Debugf("container output= nil\n")
	}
	if err != nil {
		return err
	}
	return nil
}

func (this *RKT) StopRT() error {
	return nil
}
