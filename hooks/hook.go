package hooks

type PreStart func() error

type PostStart func(containerout string) error

func SetPrestartHooks(hookpre PreStart) error {
	return hookpre()
}

func SetPostStartHooks(containerout string, hookpost PostStart) error {
	return hookpost(containerout)
}
