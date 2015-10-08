// +build v0.1.1

package specplatform

import (
	"errors"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"github.com/opencontainers/specs"
	"runtime"
)

var TestSuitePlatform manager.TestSuite = manager.TestSuite{Name: "LinuxSpec.Spec.Platform"}

func init() {
	TestSuitePlatform.AddTestCase("TestPlatformCorrect", TestPlatformCorrect)
	TestSuitePlatform.AddTestCase("TestPlatformErr", TestPlatformErr)
	manager.Manager.AddTestSuite(TestSuitePlatform)
}

func setPlatform(osValue string, archValue string) specs.LinuxSpec {
	linuxSpec := specsinit.SetLinuxspecMinimum()
	linuxSpec.Platform.OS = osValue
	linuxSpec.Platform.Arch = archValue
	return linuxSpec
}

func testPlatform(linuxSpec *specs.LinuxSpec, linuxRuntime *specs.LinuxRuntimeSpec, osValue string, archValue string) (string, error) {
	configFile := "./config.json"
	linuxSpec.Spec.Process.Args[0] = "/bin/ls"
	err := configconvert.LinuxSpecToConfig(configFile, linuxSpec)
	if err != nil {
		return manager.UNKNOWNERR, errors.New("Met unexpected LinuxSpecToConfig err, plz report to OCT project with a issuse")
	}

	rFile := "runtime.json"
	err = configconvert.LinuxRuntimeToConfig(rFile, linuxRuntime)
	if err != nil {
		return manager.UNKNOWNERR, errors.New("Met unexpected LinuxRuntimeToConfig err, plz report to OCT project with a issuse")
	}

	out, err := adaptor.StartRunc(configFile, rFile)
	if err != nil {
		if osValue != runtime.GOOS || archValue != runtime.GOARCH {
			return manager.PASSED, nil
		} else {
			return manager.FAILED, errors.New(string(out) + err.Error())
		}
	}
	if osValue == runtime.GOOS && archValue == runtime.GOARCH {
		return manager.PASSED, nil
	} else {
		return manager.FAILED, errors.New("Err: Gived err value to platform, but runc ran well")
	}
}
