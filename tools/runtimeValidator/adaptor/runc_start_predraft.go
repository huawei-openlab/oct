// +build predraft

package adaptor

import (
	"os"
	"os/exec"
)

func StartRunc(configFile string) (string, error) {
	var cmd *exec.Cmd
	if configFile == "" {
		cmd = exec.Command("runc")
	} else {
		cmd = exec.Command("runc", configFile)
	}

	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = "./"
	outPut, err := cmd.Output()
	return string(outPut), err
}
