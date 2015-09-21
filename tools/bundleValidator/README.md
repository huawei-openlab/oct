#Bundle tool

##Generate a config.json and runtime.json
Provide a template config.json and runtime.json.
`
./bundle -o config.json gc
./bundle -o runtime.json gr
`
or 
`
./bundle gc > config.json
./bundle gr > runtime.json
`

##Verify a bundle
It verifies whether a bundle is valid, with all the required files and
all the required attributes.
`
./bundle vb demo-bundle
./bundle vc demo-bundle/config.json
./bundle vr demon-bundle/runtime.json
`

#Development design
##One spec struct, one validate function
The validation work is done by .go files in the `libspec` directory.
These .go files follows the .go files in [specs](https://github.com/opencontainers/specs) closely
in order to make the validation clearly, for example, here is the 'Spec' struct:

```
opencontainers/specs/config.go

// Spec is the base configuration for the container.  It specifies platform
// independent configuration.
type Spec struct {
	// Version is the version of the specification that is supported.
	Version string `json:"version"`
	// Platform is the host information for OS and Arch.
	Platform Platform `json:"platform"`
	// Process is the container's main process.
	Process Process `json:"process"`
	// Root is the root information for the container's filesystem.
	Root Root `json:"root"`
	// Hostname is the container's host name.
	Hostname string `json:"hostname"`
	// Mounts profile configuration for adding mounts to the container's filesystem.
	MountPoints []MountPoint `json:"mounts"`
}
```

The 'Valid' function is like this:
```
libspec/config.go

func SpecValid(s specs.Spec, runtime specs.RuntimeSpec, rootfs string, msgs []string) (bool, []string) {
        valid, msgs := checkSemVer(s.Version, msgs)

        ret, msgs := PlatformValid(s.Platform, msgs)
        valid = ret && valid

        ret, msgs = ProcessValid(s.Process, msgs)
        valid = ret && valid

        ret, msgs = RootValid(s.Root, msgs)
        valid = ret && valid

        ret, msgs = StringValid("Spec.Hostname", s.Hostname, msgs)
        valid = ret && valid

        if len(s.MountPoints) > 0 {
                for index := 0; index < len(s.MountPoints); index++ {
                        ret, msgs = MountPointValid(s.MountPoints[index], msgs)
                        valid = ret && valid
                }
        }
        return valid, msgs
}
```

##Validate once, return all the erros
The return value '(bool, msgs []string)' will store all the error messages.
Correct all of them before run an OCI bundle.

```
The mountPoint sys /sys is not exist in rootfs
The mountPoint proc /proc is not exist in rootfs
The mountPoint dev /dev is not exist in rootfs
```
