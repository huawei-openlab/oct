The `oci-convert` converts:

## from aci to oci
So far we just attempt to do this.

## spec/config
Verify whether a config file is valid, with all the `required` configurations
and all the required format according to [specs](https://github.com/opencontainers/specs).

The validation work is done by .go files in the `sv` directory.
These .go files follows the .go files in [specs](https://github.com/opencontainers/specs) closely
in order to make the validation clearly:

### spec/config.go
```
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

### sv/config.go
```
```

#How To try
It is easy to use this tool, we provide a `demo rkt file` with an `manfest` file. 


```
make
./oci-convert a2i manfest
```
