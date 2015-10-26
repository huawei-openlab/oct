// +build v0.1.1

package specsinit

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/utils/configconvert"
	"testing"
)

func TestSetLinuxspecMinimum(t *testing.T) {
	ls := SetLinuxspecMinimum()
	filePath := "config.json"
	err := configconvert.LinuxSpecToConfig(filePath, &ls)
	if err != nil {
		t.Errorf("Configset TestConfigminimumset err %v", err)
	} else {
		lsn, err := configconvert.ConfigToLinuxSpec(filePath)
		if err != nil {
			t.Errorf("Configset TestConfigminimumset err %v", err)
		} else {
			if lsn.Hostname == "zenlinHost" {
				t.Log("Configset TestConfigminimumset successful!")
			} else {
				t.Error("Configset TestConfigminimumset err get wrong value from obj")
			}
		}
	}

}
