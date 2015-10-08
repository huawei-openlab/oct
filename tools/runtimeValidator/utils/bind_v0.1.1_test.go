package utils

import (
	"github.com/opencontainers/specs"
	"log"
	"os"
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
		},
	},
}

func TestBind(t *testing.T) {

	result := os.Getenv("GOPATH")
	if result == "" {
		log.Fatalf("utils.setBind error GOPATH == nil")
	}
	value := result + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/containerend"

	SetBind(&linuxRuntime, &linuxSpec)
	if linuxSpec.Mounts[1].Name == "bind" && linuxSpec.Mounts[1].Path == "/containerend" && linuxRuntime.RuntimeSpec.Mounts["bind"].Source == value {
		t.Log("TestBind SetBind successful!")
	} else {
		t.Error("TestBind SetBind err, cannot get the right value")
	}
}
