package main

import (
	"io"
	// "log"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/huawei-openlab/oct/factory"
	"github.com/huawei-openlab/oct/utils"
)

func validate(validateObj string, configArgs string) error {

	generateConfigs(validateObj, configArgs)
	prepareBundle(validateObj)

	myruntime, err := factory.CreateRuntime(Runtime)
	if err != nil {
		logrus.Printf("Create runtime %v err: %v\n", Runtime, err)
	}

	testRoot := "./bundles/" + validateObj
	if err := myruntime.StartRT(testRoot); err != nil {
		logrus.Printf("Run runtime err: %v\n", err)
		return err
	}

	return nil
}

func generateConfigs(validateObj string, configArgs string) {
	args := splitArgs(configArgs)

	logrus.Debugf("Args to the ocitools generate: ")
	for _, a := range args {
		logrus.Debugln(a)
	}
	Mutex.Lock()
	_, err := utils.ExecGenCmd(args)
	if err != nil {
		logrus.Fatalf("Generate *.json err: %v\n", err)
	}

	copy("./plugins/runtime.json-"+validateObj, "./plugins/runtime.json")
	if err != nil {
		logrus.Fatalf("copy to runtime.json-%v, %v", validateObj, err)
	}

	copy("./plugins/config.json-"+validateObj, "./plugins/config.json")
	if err != nil {
		logrus.Fatalf("copy to config.json-%v, %v", validateObj, err)
	}
	Mutex.Unlock()
}

func splitArgs(args string) []string {

	argsnew := strings.TrimSpace(args)

	argArray := strings.Split(argsnew, "--")

	lenth := len(argArray)
	resArray := make([]string, lenth-1)
	for i, arg := range argArray {
		if i == 0 || i == lenth {
			continue
		} else {
			resArray[i-1] = "--" + strings.TrimSpace(arg)
		}
	}
	return resArray
}

func prepareBundle(validateObj string) {
	// Create bundle follder
	testRoot := "./bundles/" + validateObj
	err := os.RemoveAll(testRoot)
	if err != nil {
		logrus.Fatalf("Remove bundle %v err: %v\n", validateObj, err)
	}

	err = os.Mkdir(testRoot, os.ModePerm)
	if err != nil {
		logrus.Fatalf("Mkdir bundle %v dir err: %v\n", testRoot, err)
	}

	// Create rootfs folder to bundle
	rootfs := testRoot + "/rootfs"
	err = os.Mkdir(rootfs, os.ModePerm)
	if err != nil {
		logrus.Fatalf("Mkdir rootfs for bundle %v err: %v\n", validateObj, err)
	}

	// Tar rootfs.tar.gz to rootfs
	out, err := utils.ExecCmd("", "tar", "-xf", "rootfs.tar.gz", "-C", rootfs)
	if err != nil {
		logrus.Fatalf("Tar roofs err: %v\n", out)
	}

	// Copy runtimtest from plugins to rootfs
	src := "./plugins/runtimetest"
	dRuntimeTest := rootfs + "/runtimetest"
	err = copy(dRuntimeTest, src)
	if err != nil {
		logrus.Fatalf("Copy runtimetest to rootfs err: %v\n", err)
	}
	err = os.Chmod(dRuntimeTest, os.ModePerm)
	if err != nil {
		logrus.Fatalf("Chmod runtimetest mode err: %v\n", err)
	}

	Mutex.Lock()
	// copy *.json to testroot and rootfs
	csrc := "./plugins/config.json-" + validateObj
	rsrc := "./plugins/runtime.json-" + validateObj
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
	Mutex.Unlock()
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
