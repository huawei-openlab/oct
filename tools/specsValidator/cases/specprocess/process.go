package specprocess

import (
	"github.com/huawei-openlab/oct/tools/specsValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/specsValidator/manager"
	"github.com/huawei-openlab/oct/tools/specsValidator/utils"
	"github.com/huawei-openlab/oct/tools/specsValidator/utils/configconvert"
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
			Path:     "rootfs_rootconfig",
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
	resource := result + "/src/github.com/huawei-openlab/oct/tools/specsValidator/containerend"
	utils.SetRight(resource, process.User.UID, process.User.GID)
	//linuxSpec.Spec.Mounts[0].Source = resource
	utils.SetBind(&linuxSpec, resource)
	//setBind(&linuxSpec)

	return linuxSpec
}
func testProcess(linuxspec *specs.LinuxSpec, supported bool) (string, error) {
	configFile := "./config.json"
	err := configconvert.LinuxSpecToConfig(configFile, linuxspec)
	_, err = adaptor.StartRunc(configFile)
	if err != nil {
		if supported {
			return manager.UNKNOWNERR, nil
		} else {
			return manager.PASSED, nil
		}
	}
	res := checkOut()
	if res {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, nil
	}

}

func checkResult(job string, value string) bool {
	result := utils.GetJob(job, "./containerend_out.txt")
	if strings.Contains(result, value) {
		return true
	}
	return false
}

func checkOut() bool {
	value := strconv.FormatInt(int64(linuxSpec.Spec.Process.User.UID), 10)
	job := "Uid"
	resultTag := false
	if checkResult(job, value) {
		resultTag = true
	} else {
		resultTag = false
	}

	value = strconv.FormatInt(int64(linuxSpec.Spec.Process.User.GID), 10)
	job = "Gid"
	if checkResult(job, value) {
		resultTag = true
	} else {
		resultTag = false
	}

	tmpValue := linuxSpec.Spec.Process.User.AdditionalGids
	job = "Groups"
	for _, tv := range tmpValue {
		tvs := strconv.FormatInt(int64(tv), 10)
		if checkResult(job, tvs) {
			resultTag = true
		} else {
			resultTag = false
		}
	}

	return resultTag
}
