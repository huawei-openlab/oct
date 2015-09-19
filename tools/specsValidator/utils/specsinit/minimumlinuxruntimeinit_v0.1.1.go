// +build v0.1.1

package specsinit

import (
	"github.com/opencontainers/specs"
)

func SetLinuxruntimeMinimum() specs.LinuxRuntimeSpec {
	var linuxRuntimeSpec specs.LinuxRuntimeSpec = specs.LinuxRuntimeSpec{
		RuntimeSpec: specs.RuntimeSpec{
			Mounts: map[string]specs.Mount{
				"proc": specs.Mount{
					Type:    "proc",
					Source:  "proc",
					Options: []string{""},
				},
			},
		},
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
	}
	return linuxRuntimeSpec
}
