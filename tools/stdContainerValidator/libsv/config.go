package specsValidator

import (
	"fmt"
	"github.com/opencontainers/specs"
	"os"
	"path"
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

func SpecValid(s specs.Spec, runtime specs.RuntimeSpec, rootfs string, msgs []string) (bool, []string) {
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

	ret, msgs = MountPointsValid(s.Mounts, runtime.Mounts, rootfs, msgs)
	valid = ret && valid

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
//config.md Each record in this array must have configuration in runtime config.
/ MountPoint describes a directory that may be fullfilled by a mount in the runtime.json.
type MountPoint struct {
	// Name is a unique descriptive identifier for this mount point.
	Name string `required`
	// Path specifies the path of the mount. The path and child directories MUST exist, a runtime MUST NOT create directories automatically to a mount point.
	Path string `required`
}
*/
//mps:mount points; rmps: runtime mount points
//We don't check the 'minimal mount points' here, we do it in config_linux.go
func MountPointsValid(mps []specs.MountPoint, rmps map[string]specs.Mount, rootfs string, msgs []string) (bool, []string) {
	ret := true
	valid := true
	for index := 0; index < len(mps); index++ {
		ret, msgs = MountPointValid(mps[index], rootfs, msgs)
		valid = ret && valid
		if ret == false {
			continue
		}
		if _, ok := rmps[mps[index].Name]; ok == false {
			valid = false && valid
			msgs = append(msgs, fmt.Sprintf("%s in config/mount is not exist in runtime/mount", mps[index].Name))
			continue
		}
		//Check if there were duplicated mount name
		for dIndex := index + 1; dIndex < len(mps); dIndex++ {
			if mps[index].Name == mps[dIndex].Name {
				msgs = append(msgs, fmt.Sprintf("%s in config/mount is duplicated", mps[index].Name))
				valid = false && valid
			}
		}
	}
	return valid, msgs
}

func MountPointValid(mp specs.MountPoint, rootfs string, msgs []string) (bool, []string) {
	valid, msgs := StringValid("MountPoint.Name", mp.Name, msgs)

	ret, msgs := StringValid("MountPoint.Path", mp.Path, msgs)
	valid = ret && valid

	mountPath := path.Join(rootfs, mp.Path)

	fi, err := os.Stat(mountPath)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("The mountPoint %s %s is not exist in rootfs", mp.Name, mp.Path))
		valid = ret && valid
	} else {
		if !fi.IsDir() {
			msgs = append(msgs, fmt.Sprintf("The mountPoint %s %s is not a valid directory", mp.Name, mp.Path))
			valid = ret && valid
		}
	}

	return valid, msgs
}
