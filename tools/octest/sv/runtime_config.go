package specsValidator

import (
	"github.com/opencontainers/specs"
)

/*
type RuntimeSpec struct {
	Mounts []Mount `required`
	Hooks Hooks `optional`
}
*/
func RuntimeValid(rt specs.RuntimeSpec, msgs []string) (bool, []string) {
	ret := true
	valid := true
	if len(rt.Mounts) == 0 {
		valid = false
		msgs = append(msgs, "Runtime.Mounts is missing")
	} else {
		for index := 0; index < len(rt.Mounts); index++ {
			ret, msgs = MountValid(rt.Mounts[index], msgs)
			valid = ret && valid
		}
	}

	ret, msgs = HooksValid(rt.Hooks, msgs)
	valid = ret && valid

	return valid, msgs

}

/*
type Hook struct {
	Path string   `requried`
	Args []string `required`
	Env  []string `optional`
}
*/
func HookValid(h specs.Hook, msgs []string) (bool, []string) {
	valid, msgs := StringValid("Hook.Path", h.Path, msgs)
	if len(h.Args) == 0 {
		valid = false
		msgs = append(msgs, "Hook.Args is missing")
	}

	return valid, msgs
}

/*
type Hooks struct {
	Prestart []Hook `optional`
	Poststop []Hook `optional`
}
*/
func HooksValid(hs specs.Hooks, msgs []string) (bool, []string) {
	return true, msgs
}

/*
type Mount struct {
	Type string `required`
	Source string `required`
	Destination string `required`
	Options []string `required`
}
*/

func MountValid(m specs.Mount, msgs []string) (bool, []string) {
	valid, msgs := StringValid("Mount.Type", m.Type, msgs)

	ret, msgs := StringValid("Mount.Source", m.Source, msgs)
	valid = ret && valid

	ret, msgs = StringValid("Mount.Destination", m.Destination, msgs)
	valid = ret && valid

	if len(m.Options) == 0 {
		valid = false
		msgs = append(msgs, "Mount.Options is missing")
	}

	return valid, msgs
}
