package specsValidator

import (
	"fmt"
	"github.com/opencontainers/specs"
	"strconv"
	"strings"
)

// Common functions
func checkSemVer(version string, msgs []string) (bool, []string) {
	valid, msgs := StringValid("Spec.Version", version, msgs)
	if valid {
		str := strings.Split(version, ".")
		if len(str) != 3 {
			valid = false
		} else {
			for index := 0; index < len(str); index++ {
				i, err := strconv.Atoi(str[index])
				if err != nil {
					valid = false
					break
				} else {
					if i < 0 {
						valid = false
						break
					}
				}
			}
		}
		if valid == false {
			msg := fmt.Sprintf("%s is not a valid version format, please read 'SemVer v2.0.0'", version)
			msgs = append(msgs, msg)
		}
	}
	return valid, msgs
}

func StringValid(key string, content string, msgs []string) (bool, []string) {
	valid := true
	if len(content) == 0 {
		valid = false
		msg := fmt.Sprintf("%s is missing", key)
		msgs = append(msgs, msg)
	}
	return valid, msgs
}

func SliceValid(key string, content []interface{}, msgs []string) (bool, []string) {
	valid := true
	if len(content) == 0 {
		valid = false
		msg := fmt.Sprintf("%s is missing", key)
		msgs = append(msgs, msg)
	}
	return valid, msgs
}

/*
// Spec is the base configuration for the container.  It specifies platform
// independent configuration.
type Spec struct{
	Version string `required; SemVer2.0`
	Platform Platform `required`
	Process Process `required`
	Root Root `required`
	Hostname string `optional`
	Mounts []MountPoint `optional`
}
*/

func SpecValid(s specs.Spec, msgs []string) (bool, []string) {
	valid, msgs := checkSemVer(s.Version, msgs)

	ret, msgs := PlatformValid(s.Platform, msgs)
	valid = ret && valid

	ret, msgs = ProcessValid(s.Process, msgs)
	valid = ret && valid

	ret, msgs = RootValid(s.Root, msgs)
	valid = ret && valid

	/* hostname is optional now
	ret, msgs = StringValid("Spec.Hostname", s.Hostname, msgs)
	valid = ret && valid
	*/

	if len(s.Mounts) > 0 {
		for index := 0; index < len(s.Mounts); index++ {
			ret, msgs = MountPointValid(s.Mounts[index], msgs)
			valid = ret && valid
		}
	}
	return valid, msgs
}

/*
// Process contains information to start a specific application inside the container.
type Process struct {
	Terminal bool `optional`
	User User `required`
	Args []string `required`
	Env []string `optonal`
	Cwd string `optional`
}
*/
func ProcessValid(p specs.Process, msgs []string) (bool, []string) {
	valid, msgs := UserValid(p.User, msgs)

	if len(p.Args) == 0 {
		valid = false
		msgs = append(msgs, "Process.Args is missing")
	}
	/* Cwd is optional now
	ret, msgs := StringValid("Process.Cwd", p.Cwd, msgs)
	valid = ret && valid
	*/
	return valid, msgs
}

/*
// Root contains information about the container's root filesystem on the host.
type Root struct {
	Path string `required`
	Readonly bool `optional`
}
*/
func RootValid(r specs.Root, msgs []string) (bool, []string) {
	valid, msgs := StringValid("Root.Path", r.Path, msgs)
	return valid, msgs
}

/*
// Platform specifies OS and arch information for the host system that the container
// is created for.
type Platform struct {
	OS string `required`
	Arch string `required`
}
*/

func PlatformValid(p specs.Platform, msgs []string) (bool, []string) {
	valid, msgs := StringValid("Platform.OS", p.OS, msgs)

	ret, msgs := StringValid("Platform.Arch", p.Arch, msgs)
	valid = ret && valid

	return valid, msgs
}

/*
// MountPoint describes a directory that may be fullfilled by a mount in the runtime.json.
type MountPoint struct {
	Name string `required`
	Path string `required`
}
*/
func MountPointValid(mp specs.MountPoint, msgs []string) (bool, []string) {
	valid, msgs := StringValid("MountPoint.Name", mp.Name, msgs)

	ret, msgs := StringValid("MountPoint.Path", mp.Path, msgs)
	valid = ret && valid
	return valid, msgs
}
