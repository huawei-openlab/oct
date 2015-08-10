package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/huawei-openlab/oct/engine/cmd"
	"github.com/huawei-openlab/oct/engine/setting"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	//Setting
	setting.SetConfig("conf/oct.conf")

	app := cli.NewApp()

	app.Name = setting.AppName
	app.Usage = setting.Usage
	app.Version = setting.Version
	app.Author = setting.Author
	app.Email = setting.Email

	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)

}
