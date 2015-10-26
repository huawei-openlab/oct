package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/config"
)

/*var caseConfig config.ConfigContainer

func init() {
	var err error

}*/

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
