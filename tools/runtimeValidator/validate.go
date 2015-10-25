package main

import (
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/factory"
)

var validateFlags = []cli.Flag{
	cli.StringFlag{Name: "runtime", Value: "runc", Usage: "path to the OCI runtime"},
	cli.StringFlag{Name: "validateScope", Valule: "overall", Usage: "test case to the OCI runtime"},
}

var validateCommand = cli.Command{
	Name:  "validate",
	Usage: "validate a OCI spec file",
	Flags: validateFlags,
	Action: func(context *cli.Context) {

		runtime := context.String("runtime")
		if runtime == "" {
			logrus.Fatalf("runtime path shouldn't be empty")
		}

		validateScope := context.String("validateScope")
		if validateScope == "" {
			logrus.Fatalf("validateScope shouldn't be empty")
		}

		validate(validateScope)

	},
}

func validate(validateScope string) error {

	rootfs := "./cases/" + validateScope + "input/rootfs"
	cPath := filepath.Join(rootfs, "config.json")
	rPath := filepath.Join(rootfs, "runtime.json")

	if err := factory.TestRuntime(runtime, rootfs); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Test succeeded.")
}
