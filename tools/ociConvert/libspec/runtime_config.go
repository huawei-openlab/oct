package specsConvert

import (
	"github.com/appc/spec/schema"
	"github.com/opencontainers/specs"
)

/*
type RuntimeSpec struct {
	Mounts map[string]Mount `json:"mounts"`
	Hooks Hooks `optional`
}
*/
func RuntimeSpecFrom(image schema.ImageManifest, pod schema.PodManifest, msgs []string) (specs.RuntimeSpec, []string) {
	var rt specs.RuntimeSpec

	rt.Mounts, msgs = RuntimeMountsFrom(image, pod, msgs)
	rt.Hooks, msgs = HooksFrom(image, msgs)

	return rt, msgs

}

/*
type Hook struct {
	Path string   `requried`
	Args []string `required`
	Env  []string `optional`
}
*/

/*
type Hooks struct {
	Prestart []Hook `optional`
	Poststop []Hook `optional`
}
*/
func HooksFrom(image schema.ImageManifest, msgs []string) (specs.Hooks, []string) {
	var hs specs.Hooks
	for index := 0; index < len(image.App.EventHandlers); index++ {
		var h specs.Hook
		eh := image.App.EventHandlers[index]
		if len(eh.Exec) == 0 {
			continue
		}
		for e_index := 0; e_index < len(eh.Exec); e_index++ {
			if e_index == 0 {
				h.Path = eh.Exec[0]
			} else {
				h.Args = append(h.Args, eh.Exec[e_index])
			}
		}
		switch eh.Name {
		case "pre-start":
			hs.Prestart = append(hs.Prestart, h)
			break
		case "post-stop":
			hs.Poststop = append(hs.Poststop, h)
			break
		}
	}
	return hs, msgs
}

/*
type Mount struct {
	Type string `required`
	Source string `required`
	Options []string `required`
}
*/

/* The following Volume and MountPoint are coming from AppC
type Volume struct {
	Name ACName `json:"name"`
	Kind string `json:"kind"`

	Source   string `json:"source,omitempty"`
	ReadOnly *bool  `json:"readOnly,omitempty"`
}

type MountPoint struct {
	Name     ACName `json:"name"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"readOnly,omitempty"`
}
*/
func RuntimeMountsFrom(image schema.ImageManifest, pod schema.PodManifest, msgs []string) (map[string]specs.Mount, []string) {
	mounts := make(map[string]specs.Mount)
	for index := 0; index < len(image.App.MountPoints); index++ {
		mp := image.App.MountPoints[index]
		for index_v := 0; index_v < len(pod.Volumes); index_v++ {
			volume := pod.Volumes[index_v]
			if mp.Name == volume.Name {
				var mount specs.Mount
				mount.Type = volume.Kind
				mount.Source = volume.Source
				mounts[mp.Name.String()] = mount
				break
			}
		}
	}
	return mounts, msgs
}
