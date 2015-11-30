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
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func init() {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
}

//Mutext be used in generate runtime.json & config.json
var Mutex = &sync.Mutex{}

func main() {
	app := cli.NewApp()
	app.Name = "oci-testing"
	app.Version = "0.0.1"
	app.Usage = "Utilities for OCI Testing," + "\n" +
		"\n    It is a light weight testing framework," +
		"\n    using ocitools and 3rd-party tools, " +
		"\n    managing configurable high coverage bundles as cases, " +
		"\n    supporting testing different runtimes."
	app.EnableBashCompletion = true
	app.BashComplete = cli.DefaultAppComplete
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "runtime, r",
			Value: "runc",
			Usage: "runtime to be tested, -r=runc or -r=rkt or -r=docker",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "all",
			Usage: "format and content to be ouputed, -o=all: ouput sucessful details and statics, -o=err-only: ouput failure details and statics",
		},
	}

	app.Action = func(c *cli.Context) {
		if os.Geteuid() != 0 {
			logrus.Fatal("ocitest should be run as root")
		}

		startTime := time.Now()
		runtime := c.String("runtime")
		output := c.String("output")

		wg := new(sync.WaitGroup)
		units.LoadTestUnits("./cases.conf")
		wg.Add(len(units.TestUnits))

		for _, tu := range units.TestUnits {
			go testRoutine(tu, runtime, wg)
		}

		wg.Wait()
		units.OutputResult(output)
		//logrus.Printf("Test runtime: %v, successed\n", Runtime)

		endTime := time.Now()
		dTime := endTime.Sub(startTime)
		logrus.Debugf("Cost time: %v\n", dTime.Nanoseconds())
	}

	logrus.SetLevel(logrus.InfoLevel)
	//logrus.SetLevel(logrus.DebugLevel)

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("Run App err %v\n", err)
	}
}

func testRoutine(unit *TestUnit, runtime string, wg *sync.WaitGroup) {
	logrus.Debugf("Test bundle name: %v, Test args: %v\n", unit.Name, unit.Args)
	if err := unit.SetRuntime(runtime); err != nil {
		logrus.Fatalf("Test runtime: failed to setup runtime %s , error: %v\n", runtime, err)
	} else {
		unit.Run()
	}
	wg.Done()
	return
}
