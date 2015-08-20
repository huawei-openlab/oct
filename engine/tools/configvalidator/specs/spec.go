package specs

// Spec is the base configuration for the container.  It specifies platform
// independent configuration.
type Spec struct {
	// Version is the version of the specification that is supported.
	Version string `json:"version" mandatory:"required" check:"SemVer v2.0.0"`
	// Platform is the host information for OS and Arch.
	Platform Platform `json:"platform" mandatory:"required"`
	// Process is the container's main process.
	Process Process `json:"process" mandatory:"required"`
	// Root is the root information for the container's filesystem.
	Root Root `json:"root" mandatory:"required"`
	// Hostname is the container's host name.
	Hostname string `json:"hostname" mandatory:"optional"`
	// Mounts profile configuration for adding mounts to the container's filesystem.
	Mounts []Mount `json:"mounts" mandatory:"optional"`
	// Hooks are the commands run at various lifecycle events of the container.
	Hooks Hooks `json:"hooks" mandatory:"optional"`
}

type Hooks struct {
	// Prestart is a list of hooks to be run before the container process is executed.
	// On Linux, they are run after the container namespaces are created.
	Prestart []Hook `json:"prestart"`
	// Poststop is a list of hooks to be run after the container process exits.
	Poststop []Hook `json:"poststop"`
}

// Mount specifies a mount for a container.
type Mount struct {
	// Type specifies the mount kind.
	Type string `json:"type" mandatory:"required"`
	// Source specifies the source path of the mount.  In the case of bind mounts on
	// linux based systems this would be the file on the host.
	Source string `json:"source" mandatory:"required"`
	// Destination is the path where the mount will be placed relative to the container's root.
	Destination string `json:"destination" mandatory:"required"`
	// Options are fstab style mount options.
	Options string `json:"options" mandatory":"optional"`
}

// Process contains information to start a specific application inside the container.
type Process struct {
	// Terminal creates an interactive terminal for the container.
	Terminal bool `json:"terminal" mandatory:"optional"`
	// User specifies user information for the process.
	User User `json:"user" mandatory:"required"`
	// Args specifies the binary and arguments for the application to execute.
	Args []string `json:"args" mandatory:"required"`
	// Env populates the process environment for the process.
	Env []string `json:"env" mandatory:"optional"`
	// Cwd is the current working directory for the process and must be
	// relative to the container's root.
	Cwd string `json:"cwd" mandatory:"optional"`
}

// Root contains information about the container's root filesystem on the host.
type Root struct {
	// Path is the absolute path to the container's root filesystem.
	Path string `json:"path" mandatory:"required"`
	// Readonly makes the root filesystem for the container readonly before the process is executed.
	Readonly bool `json:"readonly" mandatory:"optional"`
}

// Platform specifies OS and arch information for the host system that the container
// is created for.
type Platform struct {
	// OS is the operating system.
	OS string `json:"os" mandatory:"required"`
	// Arch is the architecture
	Arch string `json:"arch" mandatory:"required"`
}

// Hook specifies a command that is run at a particular event in the lifecycle of a container.
type Hook struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
	Env  []string `json:"env"`
}
