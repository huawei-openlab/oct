// Copyright 2015 The oct Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	sv "./libspec"
	//"fmt"
	"github.com/opencontainers/specs"
	//	"os"
)

func genConfig() (ls specs.LinuxSpec) {
	var s specs.Spec
	var plat specs.Platform
	var user specs.User
	var pro specs.Process
	var root specs.Root

	s.Version = sv.Version

	plat.OS = "linux"
	plat.Arch = "x86"
	s.Platform = plat

	user.UID = 0
	user.GID = 0
	pro.Terminal = false
	pro.User = user
	pro.Args = []string{
		"/usr/bin/ls",
		"-al",
	}
	//	pro.Env = []
	pro.Cwd = "/tmp"
	s.Process = pro

	root.Path = "rootfs"
	root.Readonly = false
	s.Root = root

	s.Hostname = "demohost"
	s.Mounts = []specs.MountPoint{
		{"sys", "/sys"},
		{"proc", "/proc"},
		{"dev", "/dev"},
		{"devpts", "/dev/pts"},
		{"devshm", "/dev/shm"},
	}

	ls.Spec = s

	var l specs.Linux
	l.Capabilities = []string{
		"CAP_AUDIT_WRITE",
		"CAP_KILL",
		"CAP_NET_BIND_SERVICE",
	}

	ls.Linux = l

	return ls
}

func genRuntime() (lrts specs.LinuxRuntimeSpec) {
	var rts specs.RuntimeSpec
	var lrt specs.LinuxRuntime

	rts.Mounts = map[string]specs.Mount{
		"sys":    specs.Mount{"sysfs", "sysfs", []string{"noexec", "nosuid", "nodev"}},
		"proc":   specs.Mount{"proc", "proc", []string{"noexec", "nosuid", "nodev"}},
		"dev":    specs.Mount{"tmpfs", "tmpfs", []string{"nosuid", "strictatime", "mode=755", "size=65536k"}},
		"devpts": specs.Mount{"devpts", "devpts", []string{"nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5"}},
		"devshm": specs.Mount{"tmpfs", "tmpfs", []string{"nosuid", "nodev"}},
	}

	lrts.RuntimeSpec = rts

	lrt.Devices = []specs.Device{
		{"/dev/random", 'c', 1, 8, "rwm", 0666, 0, 0},
		{"/dev/urandom", 'c', 1, 9, "rwm", 0666, 0, 0},
		{"/dev/null", 'c', 1, 3, "rwm", 0666, 0, 0},
		{"/dev/zero", 'c', 1, 5, "rwm", 0666, 0, 0},
		{"/dev/tty", 'c', 5, 0, "rwm", 0666, 0, 0},
		{"/dev/full", 'c', 1, 7, "rwm", 0666, 0, 0},
		{"/dev/console", 'c', 5, 1, "rwm", 0666, 0, 0},
	}

	lrt.Namespaces = []specs.Namespace{
		{"pid", ""},
		{"network", ""},
		{"mount", ""},
		{"ipc", ""},
		{"uts", ""},
		{"user", ""},
	}

	lrts.Linux = lrt

	return lrts
}
