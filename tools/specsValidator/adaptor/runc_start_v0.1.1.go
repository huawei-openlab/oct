// +build v0.1.1

package adaptor

import (
	"os"
	"os/exec"
)

func StartRunc(configFile string, runtimeFile string) (string, error) {
	var cmd *exec.Cmd
	if configFile == "" && runtimeFile == "" {
		cmd = exec.Command("runc", "start")
	} else {
		if configFile == "" {
			cmd = exec.Command("runc", "start", "-r", runtimeFile)
		} else if runtimeFile == "" {
			cmd = exec.Command("runc", "start", "-c", configFile)
		} else {
			cmd = exec.Command("runc", "start", "-c", configFile, "-r", runtimeFile)
		}
	}

	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = "./"
	outPut, err := cmd.Output()
	return string(outPut), err
}
