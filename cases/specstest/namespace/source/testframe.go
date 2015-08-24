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
	"encoding/json"
	"log"
	"reflect"
	"strings"
)

type TestCase struct {
	Name      string
	Implement reflect.Value
}

type TestSuite struct {
	Name       string
	TestCase   []*TestCase
	TestResult []string
}

const (
	PASSED     = "passed"
	FAILED     = "failed"
	UNSPPORTED = "unspported"
)

type TestResult struct {
	TestCaseName string `json:"testcasename"`
	//prepare json string
	Input interface{} `json:"input"`
	//funtion return error
	Err string `json:"error,omitempty"`
	//test result,passed,failed or unspported
	Result string `json:"result"`
}

func (this *TestResult) Set(testCaseName string, input interface{}, err error, result string) {
	this.TestCaseName = testCaseName
	this.Input = input
	if err != nil {
		this.Err = err.Error()
	} else {
		this.Err = ""
	}
	this.Result = result
}

func (this *TestResult) Marshal() string {
	js, err := json.Marshal(*this)
	if err != nil {
		log.Fatalf("Marshal error,%v\n", err)
	}
	return string(js)
}

func (this *TestSuite) AddTestCase(testCaseName string, implement interface{}) {
	for _, tc := range this.TestCase {
		if tc.Name == testCaseName || tc.Implement == reflect.ValueOf(implement) {
			log.Fatalf("Exist same testcase.")
		}
	}
	tc := new(TestCase)
	tc.Name = testCaseName
	tc.Implement = reflect.ValueOf(implement)
	this.TestCase = append(this.TestCase, tc)
}

func (this *TestSuite) Run() {

	for _, tc := range this.TestCase {
		var rt []reflect.Value
		rt = tc.Implement.Call(nil)
		str := rt[0].String()
		if str != "" {
			this.TestResult = append(this.TestResult, str)
		}
	}
}
func (this *TestSuite) GetResult() string {

	result := "{\"" + this.Name + "\":["
	for _, tr := range this.TestResult {
		result = result + tr + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "]}"
	return result
}
