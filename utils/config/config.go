package config

import (
	"bufio"
	// "fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/config"
)

var BundleMap = make(map[string]string)

/*var BundleNames = make([]string, 1)*/
var ConfigPath string
var ConfigLen int

func ReadConfig(filepath string) {
	ConfigPath = filepath
	f, err := os.Open(filepath)
	if err != nil {
		logrus.Fatalf("Open file %v error %v", filepath, err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	count := 0

	for {

		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		prefix := strings.Split(line, "=")
		caseName := strings.TrimSpace(prefix[0])
		caseArg := strings.TrimPrefix(line, caseName+"=")
		for i, arg := range splitArgs(caseArg) {
			BundleMap[caseName+strconv.FormatInt(int64(i), 10)] = arg
			count = count + 1
		}

		/*if count == 1 {
			BundleNames[0] = caseName
		} else {
			BundleNames = append(BundleNames, caseName)
		}*/
	}
	ConfigLen = count
}

func splitArgs(args string) []string {

	argArray := strings.Split(args, ";")
	resArray := make([]string, len(argArray))
	for count, arg := range argArray {
		resArray[count] = strings.TrimSpace(arg)
	}
	return resArray
}

func GetConfig(caseName string) []string {
	caseConfig, err := config.NewConfig("ini", ConfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	data := caseConfig.Strings(caseName)
	if data == nil {
		logrus.Fatalf("Get case config err.")
	}
	return data
}
