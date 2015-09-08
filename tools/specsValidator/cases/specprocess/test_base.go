package specprocess

import (
	"github.com/huawei-openlab/oct/tools/specsValidator/manager"
	"github.com/opencontainers/specs"
)

func TestBase() string {
	var process specs.Process = specs.Process{
		Terminal: false,
		User: specs.User{
			UID:            0,
			GID:            0,
			AdditionalGids: []int32{},
		},
		Args: []string{"./specprocess"},
		Env:  nil,
		Cwd:  "/containerend",
	}

	linuxspec := setProcess(process)
	newProcess := linuxspec.Spec.Process
	result, err := testProcess(&linuxspec, true)
	var testResult manager.TestResult
	testResult.Set("TestBase", newProcess, err, result)
	return testResult.Marshal()
}
