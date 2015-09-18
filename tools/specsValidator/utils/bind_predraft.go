// +build predraft

package utils

import (
	"github.com/opencontainers/specs"
)

func SetBind(linuxSpec *specs.LinuxSpec, source string) {
	linuxSpec.Mounts[0].Source = source
}
