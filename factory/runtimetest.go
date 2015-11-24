package factory

import (
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

func TestRuntime(runtime string, specDir string) error {
	logrus.Debugf("Launcing runtime")

	cmd := exec.Command(runtime, "start")
	cmd.Dir = specDir
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	logrus.Debugf("Command done")
	if string(out) != "" {
		logrus.Printf("container output=%s\n", out)
	} else {
		logrus.Debugf("container output= nil\n")
	}
	if err != nil {
		return err
	}
	return nil
}
