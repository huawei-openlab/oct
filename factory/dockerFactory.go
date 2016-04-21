package factory

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

type Docker struct {
	name string
	ID   string
}

func (this *Docker) init() {
	this.name = "docker"
	this.ID = ""
}

func (this *Docker) GetRT() string {
	return this.name
}

func (this *Docker) GetRTID() string {
	return this.ID
}

func (this *Docker) Convert(caseName string, workingDir string) (string, error) {
	var cmd *exec.Cmd
	imageName := caseName + "_docker"
	cmd = exec.Command("../plugins/oci2docker", "--debug", "--image-name", imageName, "--oci-bundle", caseName)
	cmd.Dir = workingDir
	cmd.Stdin = os.Stdin

	out, err := cmd.CombinedOutput()

	logrus.Debugf("Command done")
	if err != nil {
		return string(out), errors.New(string(out) + err.Error())
	}

	return string(out), nil
}

func (this *Docker) StartRT(specDir string) (string, error) {
	logrus.Debugf("Launcing runtime")

	caseName := filepath.Base(specDir)
	imageName := caseName + "_docker"
	casePath := filepath.Dir(specDir)

	if retStr, err := this.Convert(caseName, casePath); err != nil {
		return retStr, err
	}

	cmd := exec.Command("docker", "run", imageName)
	cmd.Dir = casePath
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	logrus.Debugf("Command done")

	id, bv, ev := this.checkResult(caseName)
	this.ID = id
	if ev != nil {
		return "", ev
	} else if !bv {
		return string(out), errors.New(string(out))
	}
	return string(out), nil
}

func (this *Docker) checkResult (caseName string) (string, bool, error) {

	//use docker ps to get uuid of docker containerr
	cmd := exec.Command("docker", "ps -a")
	cmd.Stdin = os.Stdin
	listOut, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Fatalf("docker ps err %v\n", err)
	}

	uuid, err := getUuid(string(listOut), caseName)
	if err != nil {
		return "", false, errors.New("can not get uuid of docker container" + caseName)
	}
	logrus.Debugf("uuid: %v\n", uuid)
	//TODO Based on the purpose of case to analyse the result

	return uuid, true, nil
}

func (this *Docker) StopRT(id string) error {

	cmd := exec.Command("docker", "rm -f", id)
	cmd.Stdin = os.Stdin
	_, _ = cmd.CombinedOutput()

	return nil
}
