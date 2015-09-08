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

func TestUser1000() string {
	var process specs.Process = specs.Process{
		Terminal: false,
		User: specs.User{
			UID:            1000,
			GID:            1000,
			AdditionalGids: []int32{0, 1000},
		},
		Args: []string{"./specprocess"},
		Env:  nil,
		Cwd:  "/containerend",
	}
	linuxspec := setProcess(process)
	newProcess := linuxspec.Spec.Process
	result, err := testProcess(&linuxspec, true)
	var testResult manager.TestResult
	testResult.Set("TestUser1000", newProcess, err, result)
	return testResult.Marshal()
}

func TestUser1() string {
	var process specs.Process = specs.Process{
		Terminal: false,
		User: specs.User{
			UID:            1,
			GID:            1,
			AdditionalGids: []int32{0},
		},
		Args: []string{"./specprocess"},
		Env:  nil,
		Cwd:  "/containerend",
	}
	linuxspec := setProcess(process)
	newProcess := linuxspec.Spec.Process
	result, err := testProcess(&linuxspec, true)
	var testResult manager.TestResult
	testResult.Set("TestUser1", newProcess, err, result)
	return testResult.Marshal()
}

func TestUsernil() string {
	var process specs.Process = specs.Process{
		Terminal: false,
		User: specs.User{
			UID:            0,
			GID:            0,
			AdditionalGids: nil,
		},
		Args: []string{"./specprocess"},
		Env:  nil,
		Cwd:  "/containerend",
	}
	linuxspec := setProcess(process)
	newProcess := linuxspec.Spec.Process
	result, err := testProcess(&linuxspec, true)
	var testResult manager.TestResult
	testResult.Set("TestUsernil", newProcess, err, result)
	return testResult.Marshal()
}
