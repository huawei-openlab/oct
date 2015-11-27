package main

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/huawei-openlab/oct/factory"
	"github.com/huawei-openlab/oct/utils"
	"github.com/huawei-openlab/oct/utils/config"
	"github.com/huawei-openlab/oct/utils/hooks"
)

const TestCacheDir = "./bundles/"

type TestUnit struct {
	//Case name
	Name string
	//Args is used to generate bundle
	Args string
	//Describle what does this unit test for. It is optional.
	Description string
	//Testopt is the term of OCI specs to be validate, it can be split from Args
	Testopt string

	BundleDir string
	Runtime   factory.Factory
	//success or failed
	Result string
	//when result == failed, ErrInfo containers err info, or, ErrInfo is empty string
	ErrInfo string
}

func LoadTestUnits(filename string) (units []TestUnit) {
	for key, value := range config.BundleMap {
		//TODO: config.BundleMap should support 'Description'
		unit := TestUnitNew(key, value, "")
		units = append(units, unit)
	}
	return units
}

func TestUnitNew(name string, args string, desc string) (unit TestUnit) {
	unit.Name = name
	unit.Args = args
	unit.Description = desc
	argsslice := strings.Fields(args)
	for i, arg := range argsslice {
		if strings.EqualFold(arg, "--args=./runtimetest") {
			unit.Testopt = strings.TrimPrefix(argsslice[i+1], "--args=")
		}
	}
	return unit
}

//Set runtime
func (unit *TestUnit) SetRuntime(runtime string) error {
	if r, err := factory.CreateRuntime(runtime); err != nil {
		logrus.Printf("Create runtime %v err: %v\n", runtime, err)
		return err
	} else {
		unit.Runtime = r
	}
	return nil
}

func (unit *TestUnit) Run() error {
	if unit.Runtime == nil {
		logrus.Fatalf("Set the runtime before run the test")
	}

	unit.GenerateConfigs()
	unit.PrepareBundle()

	out, err := unit.Runtime.StartRT(unit.BundleDir)
	if err != nil {
		return err
	}
	if err = unit.PostStartHooks(unit.Testopt, out); err != nil {
		return err
	}
	_ = unit.Runtime.StopRT(unit.Runtime.GetRTID())
	return nil
}

func (unit *TestUnit) PostStartHooks(testopt string, out string) error {
	var err error
	switch testopt {
	case "vna":
		err = hooks.SetPostStartHooks(out, hooks.NamespacePostStart)
	default:
	}
	return err
}

func (unit *TestUnit) PrepareBundle() {
	// Create bundle follder
	unit.BundleDir = path.Join(TestCacheDir, unit.Name)
	err := os.RemoveAll(unit.BundleDir)
	if err != nil {
		logrus.Fatalf("Remove bundle %v err: %v\n", unit.Name, err)
	}

	err = os.Mkdir(unit.BundleDir, os.ModePerm)
	if err != nil {
		logrus.Fatalf("Mkdir bundle %v dir err: %v\n", unit.BundleDir, err)
	}

	// Create rootfs folder to bundle
	rootfs := unit.BundleDir + "/rootfs"
	err = os.Mkdir(rootfs, os.ModePerm)
	if err != nil {
		logrus.Fatalf("Mkdir rootfs for bundle %v err: %v\n", unit.Name, err)
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
	csrc := "./plugins/config.json-" + unit.Name
	rsrc := "./plugins/runtime.json-" + unit.Name
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

	cdest = unit.BundleDir + "/config.json"
	rdest = unit.BundleDir + "/runtime.json"
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

func (unit *TestUnit) GenerateConfigs() {
	args := splitArgs(unit.Args)

	logrus.Debugf("Args to the ocitools generate: ")
	for _, a := range args {
		logrus.Debugln(a)
	}
	Mutex.Lock()
	_, err := utils.ExecGenCmd(args)
	if err != nil {
		logrus.Fatalf("Generate *.json err: %v\n", err)
	}

	copy("./plugins/runtime.json-"+unit.Name, "./plugins/runtime.json")
	if err != nil {
		logrus.Fatalf("copy to runtime.json-%v, %v", unit.Name, err)
	}

	copy("./plugins/config.json-"+unit.Name, "./plugins/config.json")
	if err != nil {
		logrus.Fatalf("copy to config.json-%v, %v", unit.Name, err)
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
