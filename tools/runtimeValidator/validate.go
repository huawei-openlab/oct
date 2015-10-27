package main

import (
	"io"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/factory"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
)

/*var c = make(chan int, 10)*/

func validate(validateObj string, configArgs string) error {

	generateConfigs(configArgs)
	prepareBundle(validateObj)
	testRoot := "./bundles/" + validateObj
	if err := factory.TestRuntime(Runtime, testRoot); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Test succeeded.")

	return nil
}

func generateConfigs(configArgs string) {
	args := splitArgs(configArgs)
	_, err := utils.ExecGenCmd(args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func splitArgs(args string) []string {

	argArray := strings.Split(args, "--")
	len := len(argArray) - 1
	resArray := make([]string, len)
	for count, arg := range argArray {
		if count != 0 {
			resArray[count-1] = "--" + strings.TrimSpace(arg)
		}
	}
	return resArray
}

func prepareBundle(validateObj string) {
	// Create bundle follder
	testRoot := "./bundles/" + validateObj
	err := os.RemoveAll(testRoot)
	if err != nil {
		logrus.Fatal(err)
	}

	err = os.MkdirAll(testRoot, os.ModePerm)
	if err != nil {
		logrus.Fatal(err)
	}

	// Create rootfs folder to bundle
	rootfs := testRoot + "/rootfs"
	err = os.MkdirAll(rootfs, os.ModePerm)
	if err != nil {
		logrus.Fatal(err)
	}

	// Tar rootfs.tar.gz to rootfs
	out, err := utils.ExecCmd("", "tar", "-xf", "rootfs.tar.gz", "-C", rootfs)
	if err != nil {
		logrus.Fatalf(out)
	}

	// Copy runtimtest from plugins to rootfs
	src := "./plugins/runtimetest"
	destRootfs := rootfs + "/runtimetest"
	err = copy(destRootfs, src)
	if err != nil {
		logrus.Fatal(err)
	}

	csrc := "./plugins/config.json"
	rsrc := "./plugins/runtime.json"
	cdest := rootfs + "/config.json"
	rdest := rootfs + "/runtime.json"
	err = copy(cdest, csrc)
	if err != nil {
		logrus.Fatal(err)
	}
	err = copy(rdest, rsrc)
	if err != nil {
		logrus.Fatal(err)
	}

	cdest = testRoot + "/config.json"
	rdest = testRoot + "/runtime.json"
	err = copy(cdest, csrc)
	if err != nil {
		logrus.Fatal(err)
	}
	err = copy(rdest, rsrc)
	if err != nil {
		logrus.Fatal(err)
	}
}

func copy(dst string, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}
