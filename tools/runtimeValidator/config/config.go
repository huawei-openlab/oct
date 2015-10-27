package config

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/config"
)

var CaseArray = make([]string, 1)
var ConfigPath = "cases.conf"
var ConfigLen int

func init() {
	f, err := os.Open(ConfigPath)
	if err != nil {
		logrus.Fatalf("Open file %v error %v", ConfigPath, err)
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
		if count == 0 {
			CaseArray[0] = caseName
		} else {
			CaseArray = append(CaseArray, caseName)
		}

		count = count + 1
		continue
	}
	ConfigLen = count - 1
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
