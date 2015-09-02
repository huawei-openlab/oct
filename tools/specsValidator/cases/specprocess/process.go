package specprocess

import (
	"errors"
	"github.com/huawei-openlab/oct/cases/specsValidator/adaptor"
	"github.com/huawei-openlab/oct/cases/specsValidator/manager"
	"github.com/huawei-openlab/oct/cases/specsValidator/utils/configconvert"
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

var TestSuiteProcess manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Process"}

func init() {
	TestSuiteProcess.AddTestCase("TestBase", TestBase)
	//TestSuiteProcess.AddTestCase("TestUserRoot", TestUserRoot)
	// TestSuiteProcess.AddTestCase("TestUserNoneRoot", TestUserNoneRoot)
}
