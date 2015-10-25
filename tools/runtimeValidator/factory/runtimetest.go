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
	logrus.Infof(string(out))
	if err != nil {
		return err
	}
	return nil
}
