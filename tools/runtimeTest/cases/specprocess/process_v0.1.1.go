// +build v0.1.1

package specprocess

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"strconv"
	"strings"
)

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

func setProcess(process specs.Process) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {

	linuxSpec := specsinit.SetLinuxspecMinimum()
	lr := specsinit.SetLinuxruntimeMinimum()

	//Bind containerend folder to runc container, thus we can get containerend guest programme
	linuxSpec.Spec.Process = process
	utils.SetBind(&lr, &linuxSpec)

	return linuxSpec, lr
}

func testProcessUser(linuxspec *specs.LinuxSpec, linuxruntime *specs.LinuxRuntimeSpec, supported bool) (string, error) {

	configFile := "./config.json"
	err := configconvert.LinuxSpecToConfig(configFile, linuxspec)
	if err != nil {
		return manager.UNKNOWNERR, errors.New("Got unexpected err when LinuxSpecConfig, plz report to OCT project with a issuse")
	}

	rFile := "runtime.json"
	err = configconvert.LinuxRuntimeToConfig(rFile, linuxruntime)
	if err != nil {
		return manager.UNKNOWNERR, errors.New("Got unexpected err when LinuxRuntimeToConfig, plz report to OCT project with a issuse")
	}

	output, err := adaptor.StartRunc(configFile, rFile)
	if err != nil {
		if supported {
			return manager.UNKNOWNERR, errors.New("Can not start runc, maybe runc is not support the specs with these input" + err.Error())
		} else {
			return manager.PASSED, nil
		}
	}

	res, errMsg := checkOutUser(output, *linuxspec)
	if res {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New(errMsg + "Input is not compliant with the specs or not supported by runc")
	}

}

func testProcessEnv(linuxspec *specs.LinuxSpec, linuxruntime *specs.LinuxRuntimeSpec, supported bool) (string, error) {

	configFile := "./config.json"
	err := configconvert.LinuxSpecToConfig(configFile, linuxspec)
	if err != nil {
		return manager.UNKNOWNERR, errors.New("Got unexpected err when LinuxSpecConfig, plz report to OCT project with a issuse")
	}

	rFile := "runtime.json"
	err = configconvert.LinuxRuntimeToConfig(rFile, linuxruntime)
	if err != nil {
		return manager.UNKNOWNERR, errors.New("Got unexpected err when LinuxRuntimeToConfig, plz report to OCT project with a issuse")
	}

	output, err := adaptor.StartRunc(configFile, rFile)
	if err != nil {
		if supported {
			return manager.UNKNOWNERR, errors.New("Can not start runc, maybe runc is not support the specs with these input" + string(output) + "----" + err.Error())
		} else {
			return manager.PASSED, nil
		}
	}
	value := linuxspec.Spec.Process.Env
	res := checkOutEnv(output, value)
	if res {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New("Input is not compliant with the specs or not supported by runc" + "----" + string(output))
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

func checkOutUser(output string, linuxSpec specs.LinuxSpec) (bool, string) {
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
