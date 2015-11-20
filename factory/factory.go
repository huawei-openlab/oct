package factory

import (
	"errors"
)

/*type Runtime struct {
	runtime string
}*/

type Factory interface {
	SetRT(runtime string)
	GetRT() string
	PreStart(configArgs string) error
	StartRT(specDir string) (string, error)
	PostStart(configArgs string, containerout string) error
	StopRT() error
}

func CreateRuntime(runtime string) (Factory, error) {
	switch runtime {
	case "runc":
		return new(Runc), nil
	case "rkt":
		return new(RKT), nil
	default:
		return nil, errors.New("Invalid runtime string")
	}
}
