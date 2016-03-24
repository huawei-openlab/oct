package factory

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

type Docker struct {
	name string
	ID   string
}

func (this *Docker) init() {
	this.name = ""
	this.ID = ""
}

func (this *Docker) SetRT(runtime string) {
	this.name = "docker"
}

func (this *Docker) GetRT() string {
	return "docker"
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
	imageName := appName + "_docker"
	casePath := filepath.Dir(specDir)

	if retStr, err := this.Convert(caseName, casePath); err != nil {
		return retStr, err
	}

	cmd := exec.Command("docker", "run", imageName)
	cmd.Dir = dockerPath
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	logrus.Debugf("Command done")

	id, bv, ev := checkResult(appName)
	this.ID = id
	if ev != nil {
		return "", ev
	} else if !bv {
		return string(out), errors.New(string(out))
	}
	return string(out), nil
}

func checkResult(appName string) (string, bool, error) {

	//use docker ps to get uuid of docker containerr
	cmd := exec.Command("docker", "ps -a")
	cmd.Stdin = os.Stdin
	listOut, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Fatalf("docker ps err %v\n", err)
	}

	uuid, err := getUuid(string(listOut), appName)
	if err != nil {
		return "", false, errors.New("can not get uuid of docker container" + appName)
	}
	logrus.Debugf("uuid: %v\n", uuid)
	//TODO Based on the purpose of case to analyse the result

	return uuid, true, nil
}

func getUuid(listOut string, appName string) (string, error) {

	line, err := getLine(listOut, appName)
	if err != nil {
		logrus.Debugln(err)
		return "", err
	}

	return splitUuid(line), nil
}

func splitUuid(line string) string {

	a := strings.Fields(line)
	return strings.TrimSpace(a[0])
}

func getLine(Out string, objName string) (string, error) {

	outArray := strings.Split(Out, "\n")
	flag := false
	var wantLine string
	for _, o := range outArray {
		if strings.Contains(o, objName) {
			wantLine = o
			flag = true
			break
		}
	}

	if !flag {
		return wantLine, errors.New("no line containers " + objName)
	}
	return wantLine, nil
}

func (this *Docker) StopRT(id string) error {

	cmd := exec.Command("docker", "rm -f", id)
	cmd.Stdin = os.Stdin
	_, _ = cmd.CombinedOutput()

	return nil
}
