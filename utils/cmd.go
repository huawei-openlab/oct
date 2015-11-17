package utils

import (
	// "errors"
	"io/ioutil"
	"log"
	//"os"
	"os/exec"
)

func ExecCmd(path string, arg1 string, args ...string) (string, error) {
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

func ExecGenCmd(args []string) (string, error) {
	var cmd *exec.Cmd

	argsNew := make([]string, len(args)+1)
	argsNew[0] = "generate"
	for i, a := range args {
		argsNew[i+1] = a
	}

	cmd = exec.Command("./ocitools", argsNew...)
	cmd.Dir = "./plugins"
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

func Execoci2aci(arg string) (string, error) {

	var cmd *exec.Cmd
	aciName := arg + ".aci"
	cmd = exec.Command("../plugins/oci2aci", "--debug", arg, aciName)
	cmd.Dir = "./bundles"
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
		//err = os.Rename("./plugins/"+aciName, "./bundles/"+aciName)
	}

	return retStr, err
}
