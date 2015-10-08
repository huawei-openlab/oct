// +build predraft

package specplatform

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"runtime"
)

func TestPlatformCorrect() string {
	linuxspec := setPlatform(runtime.GOOS, runtime.GOARCH)
	platform := linuxspec.Spec.Platform
	result, err := testPlatform(&linuxspec, runtime.GOOS, runtime.GOARCH)
	var testResult manager.TestResult
	testResult.Set("TestPlatformCorrect", platform, err, result)
	return testResult.Marshal()
}

func TestPlatformErr() string {
	osErr := "osErr"
	archErr := "archErr"
	linuxspec := setPlatform(osErr, archErr)
	platform := linuxspec.Spec.Platform
	result, err := testPlatform(&linuxspec, osErr, archErr)
	var testResult manager.TestResult
	testResult.Set("TestPlatformErr", platform, err, result)
	return testResult.Marshal()
}
