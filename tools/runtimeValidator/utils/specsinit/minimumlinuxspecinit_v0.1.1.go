// +build v0.1.1

package specsinit

import (
	"github.com/opencontainers/specs"
	"runtime"
)

func SetLinuxspecMinimum() specs.LinuxSpec {
	var linuxSpec specs.LinuxSpec = specs.LinuxSpec{
		Spec: specs.Spec{
			Version: "0.2.0",
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
				Args: []string{""},
			},
			Mounts: []specs.MountPoint{
				{
					Name: "proc",
					Path: "/proc",
				},
			},
		},
	}

	return linuxSpec
}
