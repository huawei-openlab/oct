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
	"os"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "oct"
	app.Usage = "Tools for OCI specs test"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:  "validate",
			Aliases: []string{"v"},
			Usage: "Validate container formats: config and layout",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "Config file to validate",
				},
				cli.StringFlag{
                                        Name:  "runtime",
                                        Usage: "Runtime file to validate",
                                },
                                cli.StringFlag{
					Name: "layout",
					Usage: "Directory layout to validate",
                                },
			},
			Action: validateProcess,
		},
	}

	app.Run(os.Args)
}

