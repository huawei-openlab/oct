// +build predraft

package specplatform

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/opencontainers/specs"
	"runtime"
)

/**
*Need mount proc and set mnt namespace when get namespace from container
*and the specs.Process.Terminal must be false when call runc in programe.
 */
var linuxSpec specs.LinuxSpec = specs.LinuxSpec{
	Spec: specs.Spec{
		Version: "pre-draft",
		Platform: specs.Platform{
			OS:   "",
			Arch: "",
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
				AdditionalGids: nil,
			},
			Args: []string{""},
		},
		Mounts: []specs.Mount{
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

var TestSuitePlatform manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Platform"}

func init() {
	TestSuitePlatform.AddTestCase("TestPlatformCorrect", TestPlatformCorrect)
	TestSuitePlatform.AddTestCase("TestPlatformErr", TestPlatformErr)
	manager.Manager.AddTestSuite(TestSuitePlatform)
}

func setPlatform(osValue string, archValue string) specs.LinuxSpec {
	linuxSpec.Platform.OS = osValue
	linuxSpec.Platform.Arch = archValue
	return linuxSpec
}

func testPlatform(linuxSpec *specs.LinuxSpec, osValue string, archValue string) (string, error) {
	configFile := "./config.json"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	linuxSpec.Spec.Process.Args[0] = "/bin/ls"
	out, err := adaptor.StartRunc(configFile)
	if err != nil {
		if osValue == runtime.GOOS && archValue == runtime.GOARCH {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New(string(out) + err.Error())
		}
	}
	if osValue == runtime.GOOS && archValue == runtime.GOARCH {
		return manager.PASSED, nil
	} else {
		return manager.UNKNOWNERR, nil
	}
}
