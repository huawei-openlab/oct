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

type RKT struct {
	name string
	ID   string
}

func (this *RKT) init() {
	this.name = "rkt"
	this.ID = ""
}

func (this *RKT) GetRT() string {
	return this.name
}

func (this *RKT) GetRTID() string {
	return this.ID
}

func (this *RKT) Convert(caseName string, workingDir string) (string, error) {
	var cmd *exec.Cmd
	aciName := caseName + ".aci"
	//set caseName to rkt appname, set rkt aciName to image name
	cmd = exec.Command("../plugins/oci2aci", "--debug", "--name", caseName, caseName, aciName)
	cmd.Dir = workingDir
	cmd.Stdin = os.Stdin

	out, err := cmd.CombinedOutput()

	logrus.Debugf("Command done")
	if err != nil {
		return string(out), errors.New(string(out) + err.Error())
	}

	return string(out), nil
}

func (this *RKT) StartRT(specDir string) (string, error) {
	logrus.Debugf("Launcing runtime")

	caseName := filepath.Base(specDir)
	aciName := caseName + ".aci"
	aciPath := filepath.Dir(specDir)

	if retStr, err := this.Convert(caseName, aciPath); err != nil {
		return retStr, err
	}

	cmd := exec.Command("rkt", "run", aciName, "--interactive", "--insecure-skip-verify", "--mds"+
		"-register=false", "--volume", "proc,kind=host,source=/bin", "--volume", "dev,kind=host,"+
		"source=/bin", "--volume", "devpts,kind=host,source=/bin", "--volume", "shm,kind=host,"+
		"source=/bin", "--volume", "mqueue,kind=host,source=/bin", "--volume", "sysfs,kind=host,"+
		"source=/bin", "--volume", "cgroup,kind=host,source=/bin", "--net=host")
	cmd.Dir = aciPath
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

func (this *RKT) checkResult (caseName string) (string, bool, error) {

	//use rkt list to get uuid of rkt container
	cmd := exec.Command("rkt", "list")
	cmd.Stdin = os.Stdin
	listOut, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Fatalf("rkt list err %v\n", err)
	}

	uuid, err := getUuid(string(listOut), caseName)
	if err != nil {
		return "", false, errors.New("can not get uuid of rkt app" + caseName)
	}
	logrus.Debugf("uuid: %v\n", uuid)

	//use rkt status to get status of app running in rkt container
	cmd = exec.Command("rkt", "status", uuid)
	statusOut, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Fatalf("rkt status err %v\n", err)
	}
	logrus.Debugf("rkt stauts %v\n,%v\n", uuid, string(statusOut))

	s, err := getAppStatus(string(statusOut), caseName)
	if s != 0 || err != nil {
		return uuid, false, err
	}

	return uuid, true, nil
}

func getAppStatus(Out string, caseName string) (int64, error) {
	line, err := getLine(Out, caseName)
	if err != nil {
		logrus.Debugln(err)
		return 1, err
	}
	a := strings.SplitAfter(line, "=")

	res, err := strconv.ParseInt(a[1], 10, 32)
	if err != nil {
		logrus.Debugln(err)
		return 1, err
	}
	return res, nil
}

func (this *RKT) StopRT(id string) error {

	cmd := exec.Command("rkt", "rm", id)
	cmd.Stdin = os.Stdin
	_, _ = cmd.CombinedOutput()

	return nil
}
