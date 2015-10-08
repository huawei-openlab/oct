// +build v0.1.1

package specplatform

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/specsinit"
	"runtime"
)

func TestPlatformCorrect() string {
	linuxspec := setPlatform(runtime.GOOS, runtime.GOARCH)
	platform := linuxspec.Spec.Platform

	lr := specsinit.SetLinuxruntimeMinimum()
	result, err := testPlatform(&linuxspec, &lr, runtime.GOOS, runtime.GOARCH)
	var testResult manager.TestResult
	testResult.Set("TestPlatformCorrect", platform, err, result)
	return testResult.Marshal()
}

func TestPlatformErr() string {
	osErr := "osErr"
	archErr := "archErr"
	linuxspec := setPlatform(osErr, archErr)
	platform := linuxspec.Spec.Platform

	lr := specsinit.SetLinuxruntimeMinimum()
	result, err := testPlatform(&linuxspec, &lr, osErr, archErr)
	var testResult manager.TestResult
	testResult.Set("TestPlatformErr", platform, err, result)
	return testResult.Marshal()
}
