package configconvert

import (
	"github.com/opencontainers/specs"
	"runtime"
	"testing"
)

var linuxSpec specs.LinuxSpec = specs.LinuxSpec{
	Spec: specs.Spec{
		Version: "0.1.0",
		Platform: specs.Platform{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
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
			Args: []string{"/bin/bash", "-c", "pwd"},
			Env:  []string{""},
			Cwd:  "/containerend",
		},
		Hostname: "zenlinHost",
		Mounts: []specs.MountPoint{
			{
				Name: "proc",
				Path: "/proc",
			},
		},
	},
}

var linuxRuntime specs.LinuxRuntimeSpec = specs.LinuxRuntimeSpec{
	Linux: specs.LinuxRuntime{
		Resources: &specs.Resources{
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
	RuntimeSpec: specs.RuntimeSpec{
		Mounts: map[string]specs.Mount{
			"proc": {
				Type:    "proc",
				Source:  "proc",
				Options: []string{""},
			},
			"bind": {
				Type:    "",
				Source:  "",
				Options: []string{""},
			},
		},
	},
}

func TestConfig_convert(t *testing.T) {
	configFile := "config.json"
	rtFile := "runtime.json"
	err := LinuxSpecToConfig(configFile, &linuxSpec)
	if err != nil {
		t.Errorf("TestConfig_convert LinuxSpecToConfig err %v", err)
	} else {
		t.Log("TestConfig_convert LinuxSpecToConfig successful!")
	}

	ls, err := ConfigToLinuxSpec(configFile)
	if err != nil {
		t.Errorf("TestConfig_convert ConfigToLinuxSpec err %v", err)
	} else {
		if ls.Hostname == "zenlinHost" {
			t.Log("TestConfig_convert ConfigToLinuxSpec successful!")
		} else {
			t.Error("TestConfig_convert ConfigToLinuxSpec err get wrong value from obj")
		}
	}

	err = LinuxRuntimeToConfig(rtFile, &linuxRuntime)
	if err != nil {
		t.Errorf("TestConfig_convert LinuxRuntimeToConfig err %v", err)
	} else {
		t.Log("TestConfig_convert LinuxRuntimeToConfig successful!")
	}

	lr, err := ConfigToLinuxRuntime(rtFile)
	if err != nil {
		t.Errorf("TestConfig_convert ConfigToLinuxRuntime err %v", err)
	} else {
		if lr.Linux.Resources.Memory.Swappiness == -1 {
			t.Log("TestConfig_convert ConfigToLinuxSpec successful!")
		} else {
			t.Error("TestConfig_convert ConfigToLinuxRuntime err get wrong value from obj")
		}
	}
}
