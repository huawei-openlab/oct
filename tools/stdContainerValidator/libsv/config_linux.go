// +build linux

package specsValidator

import (
	"fmt"
	"github.com/opencontainers/specs"
)

/*
// LinuxSpec is the full specification for linux containers.
type LinuxSpec struct {
	Spec `required`
	// Linux is platform specific configuration for linux based containers.
	Linux Linux `required`
}
*/

func LinuxSpecValid(ls specs.LinuxSpec, runtime specs.LinuxRuntimeSpec, rootfs string, msgs []string) (bool, []string) {
	var found bool

	valid, msgs := SpecValid(ls.Spec, runtime.RuntimeSpec, rootfs, msgs)

	/* The `requiredDevice` is runtime which cannot be verify by this tool
	 */
	paths := requiredPaths()
	for p_index := 0; p_index < len(paths); p_index++ {
		found = false
		for m_index := 0; m_index < len(ls.Spec.Mounts); m_index++ {
			mp := ls.Spec.Mounts[m_index]
			if paths[p_index] == mp.Path {
				found = true
			}
		}
		if found == false {
			msgs = append(msgs, fmt.Sprintf("The mount %s is missing", paths[p_index]))
			valid = found && valid
		}
	}

	ret, msgs := LinuxValid(ls.Linux, msgs)
	valid = ret && valid

	return valid, msgs
}

/*
// Linux contains platform specific configuration for linux based containers.
type Linux struct {
	// Capabilities are linux capabilities that are kept for the container.
	Capabilities []string `optional`
	// RootfsPropagation is the rootfs mount propagation mode for the container.
	RootfsPropagation string `optional`
}
*/
func LinuxValid(l specs.Linux, msgs []string) (bool, []string) {
	valid := true
	for index := 0; index < len(l.Capabilities); index++ {
		capability := "CAP_" + l.Capabilities[index]
		if capValid(capability) == false {
			msgs = append(msgs, fmt.Sprintf("%s is not valid, please `man capabilities`", l.Capabilities[index]))
			valid = false
		}

	}
	switch l.RootfsPropagation {
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

/*
// User specifies linux specific user and group information for the container's
// main process.
type User struct {
	// UID is the user id.
	UID int32 `required`
	// GID is the group id.
	GID int32 `required`
	// AdditionalGIDs are additional group ids set for the container's process.
	AdditionalGIDs []int32 `optional`
}
*/

//TODO: check if uid/gid real exist
func UserValid(u specs.User, msgs []string) (bool, []string) {
	valid := true
	if u.UID < 0 {
		valid = false
		msgs = append(msgs, "User.UID invalid")
	}
	if u.GID < 0 {
		valid = false
		msgs = append(msgs, "User.GID invalid")
	}
	return valid, msgs
}
