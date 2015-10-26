package setting

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

var (
	conf config.ConfigContainer
)

var (
	//Global
	AppName       string
	Usage         string
	Version       string
	Author        string
	Email         string
	EngineName    string
	EngienVersion string
	SpecVersion   string
)

func SetConfig(path string) error {
	var err error

	conf, err = config.NewConfig("ini", path)
	if err != nil {
		fmt.Errorf("Read %s error: %v", path, err.Error())
	}

	if appname := conf.String("appname"); appname != "" {
		AppName = appname
	} else if appname == "" {
		err = fmt.Errorf("AppName config value is null")
	}

	if usage := conf.String("usage"); usage != "" {
		Usage = usage
	} else if usage == "" {
		err = fmt.Errorf("Usage config value is null")
	}

	if version := conf.String("version"); version != "" {
		Version = version
	} else if version == "" {
		err = fmt.Errorf("Version config value is null")
	}

	if author := conf.String("author"); author != "" {
		Author = author
	} else if author == "" {
		err = fmt.Errorf("Author config value is null")
	}

	if email := conf.String("email"); email != "" {
		Email = email
	} else if email == "" {
		err = fmt.Errorf("Email config value is null")
	}

	if enginename := conf.String("enginename"); enginename != "" {
		EngineName = enginename
	} else if enginename == "" {
		err = fmt.Errorf("EngineName config value is null")
	}

	if engineversion := conf.String("enginename"); engineversion != "" {
		EngienVersion = engineversion
	} else if engineversion == "" {
		err = fmt.Errorf("EngienVersion config value is null")
	}

	if specversion := conf.String("specversion"); specversion != "" {
		SpecVersion = specversion
	} else if specversion == "" {
		err = fmt.Errorf("SpecVersion config value is null")
	}
	return err
}
