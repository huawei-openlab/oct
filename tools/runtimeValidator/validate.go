package main

import (
	"io"
	// "log"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/factory"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
)

/*var c = make(chan int, 10)*/

func validate(validateObj string, configArgs string) error {
	// logrus.Printf("validate configArgs %v\n", configArgs)

	generateConfigs(configArgs)
	prepareBundle(validateObj)
	testRoot := "./bundles/" + validateObj
	if err := factory.TestRuntime(Runtime, testRoot); err != nil {
		logrus.Fatal(err)
	}
	return nil
}

func generateConfigs(configArgs string) {
	// logrus.Printf("configArgs : %v\n", configArgs)
	args := splitArgs(configArgs)

	// logrus.Println("generateConfigs --------")
	logrus.Debugf("Args to the ocitools generate: ")
	for _, a := range args {
		logrus.Debugln(a)
	}

	_, err := utils.ExecGenCmd(args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func splitArgs(args string) []string {

	argsnew := strings.TrimSpace(args)
	// logrus.Printf("splitArgs %v\n", argsnew)

	argArray := strings.Split(argsnew, "--")

	// for _, a := range argArray {
	// 	logrus.Printf("after split %v\n", a)
	// }
	// logrus.SetLevel(logrus.WarnLevel)
	lenth := len(argArray)
	// logrus.Debugf("len : %v\n", lenth)

	// logrus.Printf("len : %v\n", lenth)
	resArray := make([]string, lenth-1)
	for i, arg := range argArray {
		// resArray[i-1] = "--" + strings.TrimSpace(arg)
		if i == 0 || i == lenth {
			continue
		} else {
			// logrus.Printf("in append %v, %v", i, arg)
			resArray[i-1] = "--" + strings.TrimSpace(arg)
			// resArray = append(resArray, "--"+arg)
		}
	}
	// logrus.Println("+++++++++++++++++++++++")

	// for _, a := range resArray {
	// 	logrus.Printf("after append  %v\n", a)
	// }

	return resArray
}

func prepareBundle(validateObj string) {
	// Create bundle follder
	testRoot := "./bundles/" + validateObj
	err := os.RemoveAll(testRoot)
	if err != nil {
		logrus.Fatal(err)
	}

	err = os.Mkdir(testRoot, os.ModePerm)
	if err != nil {
		logrus.Fatal(err)
	}

	// Create rootfs folder to bundle
	rootfs := testRoot + "/rootfs"
	err = os.Mkdir(rootfs, os.ModePerm)
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
	dRuntimeTest := rootfs + "/runtimetest"
	err = copy(dRuntimeTest, src)
	if err != nil {
		logrus.Fatal(err)
	}
	err = os.Chmod(dRuntimeTest, os.ModePerm)
	if err != nil {
		logrus.Fatal(err)
	}

	// copy *.json to testroot and rootfs
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

	csrc = "./plugins/config.json"
	rsrc = "./plugins/runtime.json"
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
