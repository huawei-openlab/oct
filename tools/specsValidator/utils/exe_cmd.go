package utils

import (
	"io/ioutil"
	"log"
	"os/exec"
)

func ExecCmdStr(arg1 string, path string, args ...string) (string, error) {
	var cmd *exec.Cmd

	argsNew := make([]string, len(args))

	for i, a := range args {
		argsNew[i] = a
	}

	cmd = exec.Command(arg1, argsNew...)
	cmd.Dir = path
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
