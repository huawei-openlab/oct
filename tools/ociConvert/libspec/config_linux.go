// +build linux

package specsConvert

import (
	"encoding/json"
	"github.com/appc/spec/schema"
	"github.com/opencontainers/specs"
	"strconv"
)

/*
// LinuxSpec is the full specification for linux containers.
type LinuxSpec struct {
	Spec `required`
	// Linux is platform specific configuration for linux based containers.
	Linux Linux `required`
}
*/

func LinuxSpecFrom(image schema.ImageManifest, msgs []string) (specs.LinuxSpec, []string) {
	var ls specs.LinuxSpec
	ls.Spec, msgs = SpecFrom(image, msgs)
	ls.Linux, msgs = LinuxFrom(image, msgs)

	return ls, msgs
}

/*
// Linux contains platform specific configuration for linux based containers.
type Linux struct {
	// Capabilities are linux capabilities that are kept for the container.
	Capabilities []string `optional`
}
*/
func LinuxFrom(image schema.ImageManifest, msgs []string) (specs.Linux, []string) {
	var l specs.Linux

	l.Capabilities, msgs = CapabilitiesFrom(image, msgs)
	return l, msgs
}

func CapabilitiesFrom(image schema.ImageManifest, msgs []string) ([]string, []string) {
	type CapSet struct {
		Set []string
	}
	var capSet CapSet
	var caps []string

	for index := 0; index < len(image.App.Isolators); index++ {
		iso := image.App.Isolators[index]
		isoName := string(iso.Name)
		switch isoName {
		case "os/linux/capabilities-retain-set":
			fallthrough
		case "os/linux/capabilities-remove-set":
			json.Unmarshal([]byte(*iso.ValueRaw), &capSet)
			for index := 0; index < len(capSet.Set); index++ {
				caps = append(caps, capSet.Set[index])
			}
		}
	}
	return caps, msgs
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

func UserFrom(image schema.ImageManifest, msgs []string) (specs.User, []string) {
	var u specs.User

	UID, err := strconv.Atoi(image.App.User)
	if err != nil {
		msgs = append(msgs, "User.UID invalid")
	} else {
		u.UID = int32(UID)
	}
	GID, err := strconv.Atoi(image.App.Group)
	if err != nil {
		msgs = append(msgs, "User.GID invalid")
	} else {
		u.GID = int32(GID)
	}

	return u, msgs
}
