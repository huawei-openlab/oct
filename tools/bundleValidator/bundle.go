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
	"./libspec"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func outputInfo(context *cli.Context, content string) {
	output := context.GlobalString("output")
	if output == "" {
		fmt.Println(content)
	} else {
		fout, err := os.Create(output)
		defer fout.Close()
		if err != nil {
			fmt.Println(output, err)
		} else {
			fout.WriteString(content)
		}
	}
}

func printErr(context *cli.Context, msgs []string) {
	output := context.GlobalString("output")
	if output == "" {
		fmt.Println(len(msgs), "errors found:")
		for index := 0; index < len(msgs); index++ {
			fmt.Println(msgs[index])
		}
	} else {
		fout, err := os.Create(output)
		defer fout.Close()
		if err != nil {
			fmt.Println(output, err)
		} else {
			for index := 0; index < len(msgs); index++ {
				fout.WriteString(msgs[index])
				fout.WriteString("\n")
			}
		}
	}
}

func parseBundle(context *cli.Context) {
	if len(context.Args()) > 0 {
		var msgs []string
		valid, msgs := specsValidator.OCTBundleValid(context.Args()[0], msgs)
		if valid {
			outputInfo(context, "Valid : config.json, runtime.json and rootfs are all accessible in the bundle")
		} else {
			printErr(context, msgs)
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
			outputInfo(context, "Valid : config.json")
		} else {
			printErr(context, msgs)
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
			outputInfo(context, "Valid : runtime.json")
		} else {
			printErr(context, msgs)
		}
	} else {
		cli.ShowCommandHelp(context, "all")
	}
}

func generateConfig(context *cli.Context) {
	ls := genConfig()
	content, _ := json.MarshalIndent(ls, "", "\t")
	outputInfo(context, string(content))
}

func generateRuntime(context *cli.Context) {
	lrt := genRuntime()
	content, _ := json.MarshalIndent(lrt, "", "\t")
	outputInfo(context, string(content))
}

func main() {
	app := cli.NewApp()
	app.Name = "Bundle Validator"
	app.Usage = "Standard Container Validator: tool to validate if a `bundle` was a standand container"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "Redirect the output to a certain file",
		},
	}
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
