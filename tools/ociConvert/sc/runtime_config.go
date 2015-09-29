package specsConvert

import (
	"github.com/appc/spec/schema"
	"github.com/opencontainers/specs"
)

/*
type RuntimeSpec struct {
	Mounts []Mount `required`
	Hooks Hooks `optional`
}
*/
func RuntimeSpecFrom(image schema.ImageManifest, msgs []string) (specs.RuntimeSpec, []string) {
	var rt specs.RuntimeSpec
	//FIXME: upstream changes!
	//	rt.Mounts, msgs = RuntimeMountsFrom(image, msgs)
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
	Destination string `required`
	Options []string `required`
}
*/

func RuntimeMountsFrom(image schema.ImageManifest, msgs []string) ([]specs.Mount, []string) {
	var mounts []specs.Mount
	//FIXME:? Where is mounts come from
	return mounts, msgs
}
