package specsConvert

import (
	"fmt"
	"github.com/appc/spec/schema"
	"github.com/opencontainers/specs"
)

/*
// Spec is the base configuration for the container.  It specifies platform
// independent configuration.
type Spec struct{
	Version string `required; SemVer2.0`
	Platform Platform `required`
	Process Process `required`
	Root Root `required`
	Hostname string `required`
	MountPoints []MountPoint `optional`
}
*/

func SpecFrom(image schema.ImageManifest, msgs []string) (specs.Spec, []string) {
	var s specs.Spec

	// pre-draft now
	s.Version = "pre-draft"
	s.Platform, msgs = PlatformFrom(image, msgs)
	s.Process, msgs = ProcessFrom(image, msgs)
	s.Root, msgs = RootFrom(image, msgs)

	s.Hostname = "missing"
	msgs = append(msgs, "hostname is not exist in aci-0.6.1")

	s.MountPoints, msgs = MountPointsFrom(image, msgs)

	return s, msgs
}

/*
// Process contains information to start a specific application inside the container.
type Process struct {
	Terminal bool `optional`
	User User `required`
	Args []string `required`
	Env []string `optonal`
	Cwd string `required`
}
*/
/* AppC
type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
*/
func EnvFrom(image schema.ImageManifest, msgs []string) ([]string, []string) {
	var env []string
	for index := 0; index < len(image.App.Environment); index++ {
		appEV := image.App.Environment[index]
		ev := fmt.Sprintf("%s=%s", appEV.Name, appEV.Value)
		env = append(env, ev)
	}
	return env, msgs
}

func ArgsFrom(image schema.ImageManifest, msgs []string) ([]string, []string) {
	return image.App.Exec, msgs
}

func ProcessFrom(image schema.ImageManifest, msgs []string) (specs.Process, []string) {
	var p specs.Process

	p.Terminal = false
	msgs = append(msgs, "terminal is not exist in aci-0.6.1")

	p.User, msgs = UserFrom(image, msgs)
	p.Args, msgs = ArgsFrom(image, msgs)
	p.Env, msgs = EnvFrom(image, msgs)
	p.Cwd = image.App.WorkingDirectory

	return p, msgs
}

/*
// Root contains information about the container's root filesystem on the host.
type Root struct {
	Path string `required`
	Readonly bool `optional`
}
*/

func RootFrom(image schema.ImageManifest, msgs []string) (specs.Root, []string) {
	var r specs.Root
	r.Path = "rootfs"
	//TODO: default to 'false'?
	r.Readonly = false
	msg := fmt.Sprintf("root readonly is not exist in aci-0.6.1")
	msgs = append(msgs, msg)
	return r, msgs
}

/*
// Platform specifies OS and arch information for the host system that the container
// is created for.
type Platform struct {
	OS string `required`
	Arch string `required`
}
*/
/* AppC
type Label struct {
	Name  ACIdentifier `json:"name"`
	Value string       `json:"value"`
}
*/

func PlatformFrom(image schema.ImageManifest, msgs []string) (specs.Platform, []string) {
	var p specs.Platform

	for index := 0; index < len(image.Labels); index++ {
		label := image.Labels[index]
		switch label.Name {
		case "os":
			p.OS = label.Value
			break
		case "arch":
			p.Arch = label.Value
			break
		default:
			break
		}
	}

	return p, msgs
}

/*
// MountPoint describes a directory that may be fullfilled by a mount in the runtime.json.
type MountPoint struct {
	Name string `required`
	Path string `required`
}
*/
func MountPointsFrom(image schema.ImageManifest, msgs []string) ([]specs.MountPoint, []string) {
	var mps []specs.MountPoint
	for index := 0; index < len(image.App.MountPoints); index++ {
		var mountPoint specs.MountPoint
		mp := image.App.MountPoints[index]
		mountPoint.Name = string(mp.Name)
		mountPoint.Path = mp.Path
		mps = append(mps, mountPoint)
	}
	msgs = append(msgs, "The mountPoint readonly att is missing during converting")
	return mps, msgs
}
