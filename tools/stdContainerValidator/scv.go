// Copyright 2015 The oct Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"./libsv"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func printErr(msgs []string) {
	fmt.Println(len(msgs), "errors found:")
	for index := 0; index < len(msgs); index++ {
		fmt.Println(msgs[index])
	}
}

func parseBundle(context *cli.Context) {
	if len(context.Args()) > 0 {
		var msgs []string
		valid := true
		os := specsValidator.OSDetect(context.Args()[0])
		if len(os) == 0 {
			valid = false
			fmt.Println("Cannot detect OS in the config.json under the bundle, or maybe miss `config.json`.")
		} else {
			if os == "linux" {
				valid, msgs = specsValidator.LinuxBundleValid(context.Args()[0], msgs)
			} else {
				valid, msgs = specsValidator.BundleValid(context.Args()[0], msgs)
			}
			if valid {
				fmt.Println("Valid : config.json, runtime.json and rootfs are all accessible in the bundle")
			} else {
				printErr(msgs)
			}
		}
	} else {
		cli.ShowCommandHelp(context, "bundle")
	}
}

func parseAll(context *cli.Context) {
	if len(context.Args()) > 0 {
		//	validateBundle(context.Args()[0])
	} else {
		cli.ShowCommandHelp(context, "all")
	}
}

// It is a cli framework.
func main() {
	app := cli.NewApp()
	app.Name = "scv"
	app.Usage = "Standard Container Validator: tool to validate if a `bundle` was a standand container"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "bundle",
			Aliases: []string{"b"},
			Usage:   "Validate if required files exist in a bundle",
			Action:  parseBundle,
		},
	}

	app.Run(os.Args)

	return
}
