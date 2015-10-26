package main

import (
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/factory"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
)

var c = make(chan int, 10)

func validate() error {

	prepare()

	rootfs := "./cases/" + validateScope + "input/rootfs"
	cPath := filepath.Join(rootfs, "config.json")
	rPath := filepath.Join(rootfs, "runtime.json")

	if err := factory.TestRuntime(runtime, rootfs); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Test succeeded.")
}

func prepare() {
	out, err := utils.ExecCmd("./", "./prepare.sh", "process")
	if err != nil {
		logrus.Fatalf(out)
	}

}
