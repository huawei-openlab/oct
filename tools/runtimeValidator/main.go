// Copyright 2015 Huawei Inc. All Rights Reserved.
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
	// "fmt"
	"os"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/config"
)

func init() {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
}

var Runtime string

func main() {
	app := cli.NewApp()
	app.Name = "oci-runtimeValidator"
	app.Version = "0.0.1"
	app.Usage = "Utilities for OCI runtime validation"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "runtime",
			Value: "runc",
			Usage: "runtime to be validated",
		},
	}
	app.Action = func(c *cli.Context) {
		Runtime = c.String("runtime")
		// for key, value := range config.BundleMap {
		// 	logrus.Println("--------------------------")
		// 	logrus.Printf("main key %v, value %v\n", key, value)
		// 	logrus.Println("--------------------------")
		// }
		for key, value := range config.BundleMap {
			logrus.Debugf("Test bundle name: %v, Test args: %v\n", key, value)
			err := validate(key, value)
			if err != nil {
				logrus.Fatal(err)
			} else {
				logrus.Printf("Test runtime: %v bundle: %v, args: %v", Runtime, key, value)
				logrus.Printf("Test Result: SUCCESS  \n")
			}
			// logrus.Debugln(err)
			// fmt.Println(err)
			time.Sleep(1 * time.Second)
		}
		//validate("process", "--args=./runtimetest --args=vp --rootfs=rootfs --terminal=false")
	}

	logrus.SetLevel(logrus.InfoLevel)
	// logrus.SetLevel(logrus.DebugLevel)

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
