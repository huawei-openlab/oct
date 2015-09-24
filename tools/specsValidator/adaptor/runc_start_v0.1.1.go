// +build v0.1.1

package adaptor

import (
	"io/ioutil"
	"log"
	"os/exec"
)

/*func StartRunc(configFile string, runtimeFile string) (string, error) {
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
}*/

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

	// cmd.stdin = os.Stdin
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("stderr err %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("stdout err %v", err)
	}

	var retStr string
	err = cmd.Start()
	if err != nil {
		retb, _ := ioutil.ReadAll(stderr)
		retStr = string(retb)
		// stdoutReader.ReadRune()
	} else {
		retb, _ := ioutil.ReadAll(stdout)
		retStr = string(retb)
	}

	return retStr, err
}
