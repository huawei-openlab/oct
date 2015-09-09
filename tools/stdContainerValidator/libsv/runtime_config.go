package specsValidator

import (
	"github.com/opencontainers/specs"
)

/*
type RuntimeSpec struct {
	Mounts map[string]Mount `json:"mounts"`
	Hooks Hooks `optional`
}
*/

/*
Minimal mount path
https://github.com/opencontainers/specs/issues/95
*/
func RuntimeSpecValid(rt specs.RuntimeSpec, msgs []string) (bool, []string) {
	ret := true
	valid := true

	for _, mount := range rt.Mounts {
		ret, msgs = MountValid(mount, msgs)
		valid = ret && valid
	}

	ret, msgs = HooksValid(rt.Hooks, msgs)
	valid = ret && valid

	return valid, msgs

}

/*
type Hook struct {
	Path string   `requried`
	Args []string `optional`
	Env  []string `optional`
}
*/
func HookValid(h specs.Hook, msgs []string) (bool, []string) {
	valid, msgs := StringValid("Hook.Path", h.Path, msgs)

	return valid, msgs
}

/*
type Hooks struct {
	Prestart []Hook `optional`
	Poststop []Hook `optional`
}
*/
func HooksValid(hs specs.Hooks, msgs []string) (bool, []string) {
	ret := true
	valid := true
	for index := 0; index < len(hs.Prestart); index++ {
		ret, msgs = HookValid(hs.Prestart[index], msgs)
		valid = ret && valid
	}
	for index := 0; index < len(hs.Poststop); index++ {
		ret, msgs = HookValid(hs.Poststop[index], msgs)
		valid = ret && valid
	}
	return valid, msgs
}

/*
type Mount struct {
	Type string `required`
	Source string `required`
	Options []string `optional`
}
*/

func MountValid(m specs.Mount, msgs []string) (bool, []string) {
	//TODO: should we check the `type`?

	valid, msgs := StringValid("Mount.Type", m.Type, msgs)

	ret, msgs := StringValid("Mount.Source", m.Source, msgs)
	valid = ret && valid

	return valid, msgs
}
