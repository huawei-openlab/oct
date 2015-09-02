// +build linux

package specsValidator

import (
	"github.com/opencontainers/specs"
)

/*
// LinuxSpec is the full specification for linux containers.
type LinuxSpec struct {
	Spec `required1
	// Linux is platform specific configuration for linux based containers.
	Linux Linux `required`
}
*/

func LinuxSpecValid(ls specs.LinuxSpec, msgs []string) (bool, []string) {
	valid, msgs := SpecValid(ls.Spec, msgs)

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
	return true, msgs
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
