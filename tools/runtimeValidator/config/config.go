package config

import (
	"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/config"
	"io"
	"os"
	"strings"
)

/*var caseConfig config.ConfigContainer

func init() {
	var err error

}*/
var CaseArray = make([]string, 1)

func GetConfig(caseName string, configPath string) []string {
	caseConfig, err := config.NewConfig("ini", configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	data := caseConfig.Strings(caseName)
	if data == nil {
		logrus.Fatalf("Get case config err.")
	}
	return data
}

func GetCaseName(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		logrus.Fatalf("Open file %v error %v", file, err)
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

		for _, c := range CaseArray {
			fmt.Println(c)
		}
		continue
	}
	return CaseArray
}
