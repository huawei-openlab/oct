package adaptor

import (
	"os/exec"
	"strings"
	"errors"
)

func GetDockerContainers() (containers []string, err error) {
	//exec shell in host machine to get docker container id
	cmd := exec.Command("/bin/sh", "-c", "docker ps -q")
	short_id, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	str := strings.TrimSuffix(string(short_id), "\n")
	if strings.TrimSpace(str) == "" {
		return nil, errors.New("can not find runing conatainers")
	}
	short_ids := strings.Split(str, "\n")
	for _, short_id := range short_ids {
		cmd = exec.Command("/bin/sh", "-c", "docker inspect -f   '{{.Id}}' "+short_id)
		full_id, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		containers = append(containers, strings.TrimSuffix(string(full_id), "\n"))
	}
	return containers, nil
}
