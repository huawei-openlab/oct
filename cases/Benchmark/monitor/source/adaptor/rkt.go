package adaptor

import (
	"errors"
	"os/exec"
	"strings"
)

func GetRktContainers() (containers []string, err error) {
	//exec shell in host machine to get rkt APP
	cmd := exec.Command("/bin/sh", "-c", "rkt list|awk '$4==\"running\" {print $2}'")
	appName, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	str := strings.TrimSuffix(string(appName), "\n")
	if strings.TrimSpace(str) == "" {
		return nil, errors.New("can not find runing conatainers")
	}
	containers = strings.Split(str, "\n")
	for i, _ := range containers {
		containers[i] = containers[i] + ".service"
	}
	return containers, nil
}
