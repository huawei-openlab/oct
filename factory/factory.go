package factory

import (
	"errors"
)

type Factory interface {
	SetRT(runtime string)
	GetRT() string
	StartRT(specDir string) (string, error)
	StopRT(id string) error
	GetRTID() string
}

func CreateRuntime(runtime string) (Factory, error) {
	switch runtime {
	case "runc":
		return new(Runc), nil
	case "rkt":
		return new(RKT), nil
	case "docker":
		return new(Docker), nil
	default:
		return nil, errors.New("Invalid runtime string")
	}
}
