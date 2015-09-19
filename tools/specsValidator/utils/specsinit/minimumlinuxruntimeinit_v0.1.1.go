// +build v0.1.1

package specsinit

import (
	"github.com/opencontainers/specs"
)

func SetLinuxruntimeMinimum() specs.LinuxRuntimeSpec {
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
	return linuxRuntime
}
