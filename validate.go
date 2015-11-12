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

/*var c = make(chan int, 10)*/

func validate(validateObj string, configArgs string) error {
	// logrus.Printf("validate configArgs %v\n", configArgs)

	generateConfigs(validateObj, configArgs)
	prepareBundle(validateObj)
	testRoot := "./bundles/" + validateObj
	if err := factory.TestRuntime(Runtime, testRoot); err != nil {
		logrus.Printf("Run runtime err: %v\n", err)
		return err
		//logrus.Fatal(err)
	}
	return nil
}

func generateConfigs(validateObj string, configArgs string) {
	// logrus.Printf("configArgs : %v\n", configArgs)
	args := splitArgs(configArgs)

	// logrus.Println("generateConfigs --------")
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
