package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/opencontainers/specs"
)

func RktParaGenVol(configpath string, runtimepath string) (string, error) {
	var spec *specs.LinuxSpec
	var rspec *specs.LinuxRuntimeSpec

	cPath := configpath
	cf, err := os.Open(cPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Errorf("config.json not found")
			return "error", fmt.Errorf("config.json not found")
		}
	}
	defer cf.Close()

	rPath := runtimepath
	rf, err := os.Open(rPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Errorf("runtime.json not found")
			return "error", fmt.Errorf("runtime.json not found")
		}
	}
	defer rf.Close()

	if err = json.NewDecoder(cf).Decode(&spec); err != nil {
		return "error", err
	}
	if err = json.NewDecoder(rf).Decode(&rspec); err != nil {
		return "error", err
	}
	outstring := ""
	for _, mnt := range spec.Mounts {
		outstring = outstring + " --volume " + mnt.Name + ",kind=host,source=" + rspec.Mounts[mnt.Name].Source
	}
	return outstring, nil
}
