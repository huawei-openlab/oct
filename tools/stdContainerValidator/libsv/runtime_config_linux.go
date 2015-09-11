package specsValidator

import (
	"fmt"
	"github.com/opencontainers/specs"
)

/*
// LinuxSpec is the full specification for linux containers.
type LinuxRuntimeSpec struct {
	RuntimeSpec
	// Linux is platform specific configuration for linux based containers.
	Linux LinuxRuntime `json:"linux"`
}
*/

func LinuxRuntimeSpecValid(lrs specs.LinuxRuntimeSpec, msgs []string) (bool, []string) {
	valid, msgs := RuntimeSpecValid(lrs.RuntimeSpec, msgs)
	ret, msgs := LinuxRuntimeValid(lrs.Linux, msgs)
	valid = ret && valid
	return valid, msgs
}

/*
type LinuxRuntime struct {
	// UIDMapping specifies user mappings for supporting user namespaces on linux.
	UIDMappings []IDMapping `optional`
	// UIDMapping specifies group mappings for supporting user namespaces on linux.
	GIDMappings []IDMapping `optional`
	// Rlimits specifies rlimit options to apply to the container's process.
	Rlimits []Rlimit `optional`
	// Sysctl are a set of key value pairs that are set for the container on start
	Sysctl map[string]string `json:"sysctl"`
	// Resources contain cgroup information for handling resource constraints
	// for the container
	Resources Resources `json:"resources"`
	// Namespaces contains the namespaces that are created and/or joined by the container
	Namespaces []Namespace `optional`
	// Devices are a list of device nodes that are created and enabled for the container
	Devices []Device `json:"devices"`
	// ApparmorProfile specified the apparmor profile for the container.
	ApparmorProfile string `json:"apparmorProfile"`
	// SelinuxProcessLabel specifies the selinux context that the container process is run as.
	SelinuxProcessLabel string `json:"selinuxProcessLabel"`
	// Seccomp specifies the seccomp security settings for the container.
	Seccomp Seccomp `json:"seccomp"`
	// RootfsPropagation is the rootfs mount propagation mode for the container
	RootfsPropagation string `json:"rootfsPropagation"`
}
*/

func LinuxRuntimeValid(lr specs.LinuxRuntime, msgs []string) (bool, []string) {
	ret := true
	valid := true
	if len(lr.UIDMappings)+len(lr.GIDMappings) > 5 {
		valid = false
		msgs = append(msgs, "The UID/GID mapping is limited to 5")
	}

	for index := 0; index < len(lr.Rlimits); index++ {
		ret, msgs = RlimitValid(lr.Rlimits[index], msgs)
		valid = ret && valid
	}

	ret, msgs = ResourcesValid(lr.Resources, msgs)
	valid = ret && valid

	for index := 0; index < len(lr.Namespaces); index++ {
		ret, msgs = NamespaceValid(lr.Namespaces[index], msgs)
		valid = ret && valid
	}

	//minimum devices
	devices := requiredDevices()
	for index := 0; index < len(devices); index++ {
		found := false
		for dIndex := 0; dIndex < len(lr.Devices); dIndex++ {
			if lr.Devices[dIndex].Path == devices[index] {
				found = true
				break
			}
		}
		if found == false {
			msgs = append(msgs, fmt.Sprintf("The required device %s is missing", devices[index]))
			valid = found && valid
		}
	}

	for index := 0; index < len(lr.Devices); index++ {
		ret, msgs = DeviceValid(lr.Devices[index], msgs)
		valid = ret && valid
	}

	switch lr.RootfsPropagation {
	case "":
	case "slave":
	case "private":
	case "shared":
		break
	default:
		valid = false
		msgs = append(msgs, "RootfsPropagation should limited to 'slave', 'private', or 'shared'")
		break

	}

	return valid, msgs
}

/* Namespace is the configuration for a linux namespace.
type Namespace struct {
	// Type is the type of Linux namespace
	Type string `required`
	// Path is a path to an existing namespace persisted on disk that can be joined
	// and is of the same type
	Path string `optional`
}
*/
func NamespaceValid(ns specs.Namespace, msgs []string) (bool, []string) {
	valid := true
	switch ns.Type {
	case "":
		valid = false
		msgs = append(msgs, "The type of the namespace should not be empty")
		break
	case "pid":
	case "network":
	case "mount":
	case "ipc":
	case "uts":
	case "user":
		break
	default:
		valid = false
		msgs = append(msgs, "The type of the namespace should limited to 'pid/network/mount/ipc/nts/user'")
		break
	}
	return valid, msgs
}

/*
// IDMapping specifies UID/GID mappings
type IDMapping struct {
	// HostID is the UID/GID of the host user or group
	HostID int32 `json:"hostID"`
	// ContainerID is the UID/GID of the container's user or group
	ContainerID int32 `json:"containerID"`
	// Size is the length of the range of IDs mapped between the two namespaces
	Size int32 `json:"size"`
}
*/

func IDMappingValid(idm specs.IDMapping, msgs []string) (bool, []string) {
	//TODO?
	return true, msgs
}

/*
// Rlimit type and restrictions
type Rlimit struct {
	// Type of the rlimit to set
	Type int `json:"type"`
	// Hard is the hard limit for the specified type
	Hard uint64 `json:"hard"`
	// Soft is the soft limit for the specified type
	Soft uint64 `json:"soft"`
}
*/

func RlimitValid(r specs.Rlimit, msgs []string) (bool, []string) {
	if rlimitValid(r.Type) {
		msgs = append(msgs, "Rlimit is invalid")
		return false, msgs
	}
	return true, msgs
}

/*
// HugepageLimit structure corresponds to limiting kernel hugepages
type HugepageLimit struct {
	Pagesize string `json:"pageSize"`
	Limit    int    `json:"limit"`
}

// InterfacePriority for network interfaces
type InterfacePriority struct {
	// Name is the name of the network interface
	Name string `json:"name"`
	// Priority for the interface
	Priority int64 `json:"priority"`
}

// BlockIO for Linux cgroup 'blockio' resource management
type BlockIO struct {
	// Specifies per cgroup weight, range is from 10 to 1000
	Weight int64 `json:"blkioWeight"`
	// Weight per cgroup per device, can override BlkioWeight
	WeightDevice string `json:"blkioWeightDevice"`
	// IO read rate limit per cgroup per device, bytes per second
	ThrottleReadBpsDevice string `json:"blkioThrottleReadBpsDevice"`
	// IO write rate limit per cgroup per divice, bytes per second
	ThrottleWriteBpsDevice string `json:"blkioThrottleWriteBpsDevice"`
	// IO read rate limit per cgroup per device, IO per second
	ThrottleReadIOpsDevice string `json:"blkioThrottleReadIopsDevice"`
	// IO write rate limit per cgroup per device, IO per second
	ThrottleWriteIOpsDevice string `json:"blkioThrottleWriteIopsDevice"`
}

// Memory for Linux cgroup 'memory' resource management
type Memory struct {
	// Memory limit (in bytes)
	Limit int64 `json:"limit"`
	// Memory reservation or soft_limit (in bytes)
	Reservation int64 `json:"reservation"`
	// Total memory usage (memory + swap); set `-1' to disable swap
	Swap int64 `json:"swap"`
	// Kernel memory limit (in bytes)
	Kernel int64 `json:"kernel"`
	// How aggressive the kernel will swap memory pages. Range from 0 to 100. Set -1 to use system default
	Swappiness int64 `json:"swappiness"`
}

// CPU for Linux cgroup 'cpu' resource management
type CPU struct {
	// CPU shares (relative weight vs. other cgroups with cpu shares)
	Shares int64 `json:"shares"`
	// CPU hardcap limit (in usecs). Allowed cpu time in a given period
	Quota int64 `json:"quota"`
	// CPU period to be used for hardcapping (in usecs). 0 to use system default
	Period int64 `json:"period"`
	// How many time CPU will use in realtime scheduling (in usecs)
	RealtimeRuntime int64 `json:"realtimeRuntime"`
	// CPU period to be used for realtime scheduling (in usecs)
	RealtimePeriod int64 `json:"realtimePeriod"`
	// CPU to use within the cpuset
	Cpus string `json:"cpus"`
	// MEM to use within the cpuset
	Mems string `json:"mems"`
}

// Network identification and priority configuration
type Network struct {
	// Set class identifier for container's network packets
	ClassID string `json:"classId"`
	// Set priority of network traffic for container
	Priorities []InterfacePriority `json:"priorities"`
}

// Resources has container runtime resource constraints
type Resources struct {
	// DisableOOMKiller disables the OOM killer for out of memory conditions
	DisableOOMKiller bool `json:"disableOOMKiller"`
	// Memory restriction configuration
	Memory Memory `json:"memory"`
	// CPU resource restriction configuration
	CPU CPU `json:"cpu"`
	// BlockIO restriction configuration
	BlockIO BlockIO `json:"blockIO"`
	// Hugetlb limit (in bytes)
	HugepageLimits []HugepageLimit `json:"hugepageLimits"`
	// Network restriction configuration
	Network Network `json:"network"`
}
*/

func ResourcesValid(r *specs.Resources, msgs []string) (bool, []string) {
	if r == nil {
		return true, msgs
	}
	return true, msgs
}

/*
type Device struct {
	// Path to the device.
	Path string `json:"path"`
	// Device type, block, char, etc.
	Type rune `json:"type"`
	// Major is the device's major number.
	Major int64 `json:"major"`
	// Minor is the device's minor number.
	Minor int64 `json:"minor"`
	// Cgroup permissions format, rwm.
	Permissions string `json:"permissions"`
	// FileMode permission bits for the device.
	FileMode os.FileMode `json:"fileMode"`
	// UID of the device.
	UID uint32 `json:"uid"`
	// GID of the device.
	GID uint32 `json:"gid"`
}
*/

func DeviceValid(d specs.Device, msgs []string) (bool, []string) {
	valid, msgs := StringValid("Device.Path", d.Path, msgs)
	if valid == false {
		return valid, msgs
	}

	switch d.Type {
	case 'b':
	case 'c':
	case 'u':
		if d.Major <= 0 {
			msgs = append(msgs, fmt.Sprintf("Device %s type is `b/c/u`, please set the major number", d.Path))
			valid = false && valid
		}
		if d.Minor <= 0 {
			msgs = append(msgs, fmt.Sprintf("Device %s type is `b/c/u`, please set the minor number", d.Path))
			valid = false && valid
		}
		break
	case 'p':
		if d.Major > 0 || d.Minor > 0 {
			msgs = append(msgs, fmt.Sprintf("Device %s type is `p`, no need to set major/minor number", d.Path))
			valid = false && valid
		}
		break
	default:
		msgs = append(msgs, fmt.Sprintf("Device %s type should limited to `b/c/u/p`", d.Path))
		valid = false && valid
		break
	}
	//TODO, check permission/filemode

	if d.UID <= 0 {
		msgs = append(msgs, fmt.Sprintf("Device %s UID %d is invalid`", d.Path, d.UID))
		valid = false && valid
	}
	if d.GID <= 0 {
		msgs = append(msgs, fmt.Sprintf("Device %s GID %d is invalid`", d.Path, d.GID))
		valid = false && valid
	}
	return true, msgs
}

/*
// Seccomp represents syscall restrictions
type Seccomp struct {
	DefaultAction Action     `json:"defaultAction"`
	Syscalls      []*Syscall `json:"syscalls"`
}

// Action taken upon Seccomp rule match
type Action string

// Operator used to match syscall arguments in Seccomp
type Operator string

// Arg used for matching specific syscall arguments in Seccomp
type Arg struct {
	Index    uint     `json:"index"`
	Value    uint64   `json:"value"`
	ValueTwo uint64   `json:"valueTwo"`
	Op       Operator `json:"op"`
}

// Syscall is used to match a syscall in Seccomp
type Syscall struct {
	Name   string `json:"name"`
	Action Action `json:"action"`
	Args   []*Arg `json:"args"`
}
*/
