package factory

import (
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

func TestRuntime(runtime string, specDir string) error {
	logrus.Infof("Launcing runtime")

	cmd := exec.Command(runtime, "start")
	cmd.Dir = specDir
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	logrus.Infof("Command done")
	if string(out) != "" {
		logrus.Infof("container output=%s", out)
	} else {
		logrus.Infof("container output= nil")
	}
	if err != nil {
		return err
	}
	return nil
}
