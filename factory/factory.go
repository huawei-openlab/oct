package factory

import (
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
)

type Factory interface {
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

func getUuid(listOut string, caseName string) (string, error) {

	line, err := getLine(listOut, caseName)
	if err != nil {
		logrus.Debugln(err)
		return "", err
	}

	return splitUuid(line), nil
}

func splitUuid(line string) string {

	a := strings.Fields(line)
	return strings.TrimSpace(a[0])
}

func getLine(Out string, objName string) (string, error) {

	outArray := strings.Split(Out, "\n")
	flag := false
	var wantLine string
	for _, o := range outArray {
		if strings.Contains(o, objName) {
			wantLine = o
			flag = true
			break
		}
	}

	if !flag {
		return wantLine, errors.New("no line containers " + objName)
	}
	return wantLine, nil
}
