// +build predraft

package specprocess

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/opencontainers/specs"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/**
*Need mount proc and set mnt namespace when get namespace from container
*and the specs.Process.Terminal must be false when call runc in programe.
 */
var linuxSpec specs.LinuxSpec = specs.LinuxSpec{
	Spec: specs.Spec{
		Version: "pre-draft",
		Platform: specs.Platform{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		Root: specs.Root{
			Path:     "rootfs",
			Readonly: true,
		},
		Process: specs.Process{
			Terminal: false,
			User: specs.User{
				UID:            0,
				GID:            0,
				AdditionalGids: []int32{1},
			},
			Args: []string{""},
			Env:  []string{""},
			Cwd:  "",
		},
		Mounts: []specs.Mount{
			{
				Type:        "bind",
				Source:      "",
				Destination: "/containerend",
				Options:     "bind",
			},
			{
				Type:        "proc",
				Source:      "proc",
				Destination: "/proc",
				Options:     "",
			},
		},
	},
	Linux: specs.Linux{
		Resources: specs.Resources{
			Memory: specs.Memory{
				Swappiness: -1,
			},
		},
		Namespaces: []specs.Namespace{
			{
				Type: "mount",
				Path: "",
			},
		},
	},
}

var TestSuiteProcess manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Process"}

func init() {
	TestSuiteProcess.AddTestCase("TestBase", TestBase)
	TestSuiteProcess.AddTestCase("TestUser1000", TestUser1000)
	TestSuiteProcess.AddTestCase("TestUser1", TestUser1)
	TestSuiteProcess.AddTestCase("TestUsernil", TestUsernil)
	TestSuiteProcess.AddTestCase("TestEnv", TestEnv)
	TestSuiteProcess.AddTestCase("TestEnvNilFalse", TestEnvNilFalse)
	TestSuiteProcess.AddTestCase("TestEnvNilTrue", TestEnvNilTrue)
	manager.Manager.AddTestSuite(TestSuiteProcess)
	//TestSuiteProcess.AddTestCase("TestUserRoot", TestUserRoot)
	// TestSuiteProcess.AddTestCase("TestUserNoneRoot", TestUserNoneRoot)
}

func setProcess(process specs.Process) specs.LinuxSpec {
	linuxSpec.Spec.Process = process
	//linuxSpec.Spec.Process.Args = append(linuxSpec.Spec.Process.Args, "/specprocess")
	//linuxSpec.Spec.Process.Args[0] = "./specprocess"

	result := os.Getenv("GOPATH")
	if result == "" {
		log.Fatalf("utils.setBind error GOPATH == nil")
	}
	resource := result + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/containerend"
	utils.SetRight(resource, process.User.UID, process.User.GID)
	//linuxSpec.Spec.Mounts[0].Source = resource
	utils.SetBind(&linuxSpec, resource)

	return linuxSpec
}
func testProcessUser(linuxspec *specs.LinuxSpec, supported bool) (string, error) {
	configFile := "./config.json"
	err := configconvert.LinuxSpecToConfig(configFile, linuxspec)
	output, err := adaptor.StartRunc(configFile)
	if err != nil {
		if supported {
			return manager.UNKNOWNERR, errors.New("Can not start runc, maybe runc is not support the specs with this config.json" + err.Error())
		} else {
			return manager.PASSED, nil
		}
	}
	res, errMsg := checkOutUser(output)
	if res {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New(errMsg + " in runtime is not compliant with the specs")
	}

}

func testProcessEnv(linuxspec *specs.LinuxSpec, supported bool) (string, error) {
	configFile := "./config.json"
	err := configconvert.LinuxSpecToConfig(configFile, linuxspec)
	output, err := adaptor.StartRunc(configFile)
	if err != nil {
		if supported {
			return manager.UNKNOWNERR, errors.New("Can not start runc, maybe runc is not support the specs with this config.json----" + string(output) + "----" + err.Error())
		} else {
			return manager.PASSED, nil
		}
	}
	value := linuxSpec.Spec.Process.Env
	res := checkOutEnv(output, value)
	if res {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New("Env in runtime is not compliant with the specs" + "----" + string(output))
	}
}

func checkOutEnv(output string, value []string) bool {
	//fmt.Println(output)
	var rt *[]bool = new([]bool)
	for _, va := range value {
		if strings.Contains(output, va) {
			if value != nil {
				*rt = append(*rt, true)
			} else {
				*rt = append(*rt, false)
			}

		} else {
			if value != nil {
				*rt = append(*rt, false)
			} else {
				*rt = append(*rt, true)
			}
		}
	}
	var tmp bool
	for i, r := range *rt {
		if i == 0 {
			tmp = r
		} else {
			tmp = tmp && r
		}
	}
	return tmp

}

func getJob(job string, output string) string {
	as := strings.Split(output, "\n")
	for _, s := range as {
		if strings.Contains(s, job) {
			return s
		}
	}
	return ""
}

func checkResult(job string, value string, output string) bool {
	result := getJob(job, output)
	if strings.Contains(result, value) {
		return true
	}
	return false
}

func checkOutUser(output string) (bool, string) {
	value := strconv.FormatInt(int64(linuxSpec.Spec.Process.User.UID), 10)
	job := "Uid"
	var rt1, rt2, rt3 bool
	if checkResult(job, value, output) {
		rt1 = true
	} else {
		rt1 = false
	}

	value = strconv.FormatInt(int64(linuxSpec.Spec.Process.User.GID), 10)
	job = "Gid"
	if checkResult(job, value, output) {
		rt2 = true
	} else {
		rt2 = false
	}

	tmpValue := linuxSpec.Spec.Process.User.AdditionalGids
	job = "Groups"

	if tmpValue != nil {
		for _, tv := range tmpValue {
			tvs := strconv.FormatInt(int64(tv), 10)
			if checkResult(job, tvs, output) {
				rt3 = true
			} else {
				rt3 = false
			}
		}
	} else {
		result := getJob(job, output)
		gids := strings.SplitAfter(result, ":")
		gid := strings.TrimSpace(gids[1])
		if gid == "" {
			rt3 = true
		} else {
			rt3 = false
		}

	}
	var errPart *[]string = new([]string)
	if !rt1 {
		*errPart = append(*errPart, "UID")
	}

	if !rt2 {
		*errPart = append(*errPart, "GID")
	}

	if !rt3 {
		*errPart = append(*errPart, "AddtionalGids")
	}

	rt := rt1 && rt2 && rt3
	if rt {
		*errPart = append(*errPart, "")
	}
	var rs string
	if !rt {
		for i, ep := range *errPart {
			if i == 0 {
				rs = ep
			} else {
				rs = rs + " | " + ep
			}
		}
	}
	return rt, rs
}
