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
	"github.com/codegangsta/cli"
	"os"
)

func convertProcess(context *cli.Context) {
	if len(context.Args()) > 1 {
		convertRocketFile(context.Args()[0], context.Args()[1])
	}

}

// It is a cli framework.
func main() {
	app := cli.NewApp()
	app.Name = "oci-convert"
	app.Usage = "Tools for convert between different container bundle/image"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:    "aci2oci",
			Aliases: []string{"a2o"},
			Usage:   "convert container formats from aci to oci: config and bundle",

			Action: convertProcess,
		},
	}

	app.Run(os.Args)

	return
}
