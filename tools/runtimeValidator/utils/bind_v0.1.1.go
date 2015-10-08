// +build v0.1.1

package utils

import (
	"github.com/opencontainers/specs"
	"log"
	"os"
)

func SetBind(linuxRuntime *specs.LinuxRuntimeSpec, linuxSpec *specs.LinuxSpec) {

	//testtoolfolder := specs.Mount{"bind", resource, "/testtool", "bind"}
	result := os.Getenv("GOPATH")
	if result == "" {
		log.Fatalf("utils.setBind error GOPATH == nil")
	}
	source := result + "/src/github.com/huawei-openlab/oct/tools/runtimeValidator/containerend"
	mountpoint := specs.MountPoint{"bind", "/containerend"}
	linuxSpec.Mounts = append(linuxSpec.Mounts, mountpoint)
	linuxRuntime.Mounts["bind"] = specs.Mount{"bind", source, []string{"bind"}}

	SetRight(source, linuxSpec.Process.User.UID, linuxSpec.Process.User.GID)
}
