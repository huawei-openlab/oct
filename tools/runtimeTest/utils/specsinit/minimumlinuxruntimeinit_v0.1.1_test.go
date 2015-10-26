// +build v0.1.1

package specsinit

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"testing"
)

func TestSetLinuxruntimeMinimum(t *testing.T) {
	lr := SetLinuxruntimeMinimum()
	filePath := "runtime.json"
	err := configconvert.LinuxRuntimeToConfig(filePath, &lr)
	if err != nil {
		t.Errorf("specsinit SetLinuxruntimeMinimum err %v", err)
	} else {
		lrn, err := configconvert.ConfigToLinuxRuntime(filePath)
		if err != nil {
			t.Errorf("specsinit SetLinuxruntimeMinimum err %v", err)
		} else {
			if lrn.Linux.Resources.Memory.Swappiness == -1 {
				t.Log("specsinit SetLinuxruntimeMinimum successful!")
			} else {
				t.Error("specsinit SetLinuxruntimeMinimum err get wrong value from obj")
			}
		}
	}

}
