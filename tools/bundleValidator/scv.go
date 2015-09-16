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
	"encoding/json"
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
		valid, msgs := specsValidator.OCTBundleValid(context.Args()[0], msgs)
		if valid {
			fmt.Println("Valid : config.json, runtime.json and rootfs are all accessible in the bundle")
		} else {
			printErr(msgs)
		}
	} else {
		cli.ShowCommandHelp(context, "bundle")
	}
}

func parseConfig(context *cli.Context) {
	if len(context.Args()) > 0 {
		var msgs []string
		valid, msgs := specsValidator.OCTConfigValid(context.Args()[0], msgs)
		if valid {
			fmt.Println("Valid : config.json")
		} else {
			printErr(msgs)
		}
	} else {
		cli.ShowCommandHelp(context, "bundle")
	}
}

func parseRuntime(context *cli.Context) {
	if len(context.Args()) > 0 {
		var msgs []string
		var os string
		if len(context.Args()) > 1 {
			os = context.Args()[1]
		}
		valid, msgs := specsValidator.OCTRuntimeValid(context.Args()[0], os, msgs)
		if valid {
			fmt.Println("Valid : runtime.json")
		} else {
			printErr(msgs)
		}
	} else {
		cli.ShowCommandHelp(context, "all")
	}
}

func generateConfig(context *cli.Context) {
	ls := genConfig()
	content, _ := json.Marshal(ls)
	fmt.Println(string(content))
}

func generateRuntime(context *cli.Context) {
	lrt := genRuntime()
	content, _ := json.Marshal(lrt)
	fmt.Println(string(content))
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
			Usage:   "Validate all the config.json, runtime.json and files in the rootfs",
			Action:  parseBundle,
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Validate the config.json only",
			Action:  parseConfig,
		},
		{
			Name:    "runtime",
			Aliases: []string{"r"},
			Usage:   "Validate the runtime.json only, runtime + arch, default to 'linux'",
			Action:  parseRuntime,
		},
		{
			Name:    "genconfig",
			Aliases: []string{"gc"},
			Usage:   "Generate a demo config.json",
			Action:  generateConfig,
		},
		{
			Name:    "genruntime",
			Aliases: []string{"gr"},
			Usage:   "Generate a demo runtime.json",
			Action:  generateRuntime,
		},
	}

	app.Run(os.Args)

	return
}
