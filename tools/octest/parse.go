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
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
)

const (
	// Path to config file inside the bundle
	ConfigFile  = "config.json"
	RuntimeFile = "runtime.json"
	// Path to rootfs directory inside the bundle
	RootfsDir = "rootfs"
)

var (
	ErrNoRootFS = errors.New("no rootfs found in bundle")
	ErrNoConfig = errors.New("no config json file found in bundle")
	ErrNoRun    = errors.New("no runtime json file found in bundle")
)

func validateProcess(context *cli.Context) {
	//parse --config, --runtime, --bundle option
	if args := context.String("config"); len(args) != 0 {
		//validate config.json
		validateConfigFile(args)
	} else if args := context.String("runtime"); len(args) != 0 {
		//validate runtime.json
		validateRuntime(args)
	} else if args := context.String("bundle"); len(args) != 0 {
		//validate bundle
		err := validateBundle(args)
		if err != nil {
			fmt.Printf("%s: invalid image bundle: %v\n", args, err)
		} else {
			fmt.Printf("%s: valid image bundle\n", args)
		}

	} else {
		cli.ShowCommandHelp(context, "validate")
		return
	}
}
