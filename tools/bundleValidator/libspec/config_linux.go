// +build linux

package specsValidator

import (
	"fmt"
	"github.com/opencontainers/specs"
)

/*
type LinuxSpec struct {
	Spec `required`
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
type Linux struct {
	Capabilities []string `optional`
}
*/
func LinuxValid(l specs.Linux, msgs []string) (bool, []string) {
	valid := true
	for index := 0; index < len(l.Capabilities); index++ {
		capability := l.Capabilities[index]
		if capValid(capability) == false {
			msgs = append(msgs, fmt.Sprintf("%s is not valid, please `man capabilities`", l.Capabilities[index]))
			valid = false
		}

	}

	return valid, msgs
}

/*
type User struct {
	UID int32 `required`
	GID int32 `required`
	AdditionalGIDs []int32 `optional`
}
*/

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
