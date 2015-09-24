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
package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

//Top level manger to organize the testsuite
type TestManager struct {
	TestSuite []*TestSuite
}

type TestCase struct {
	// Testcase name, assignment in AddTestCase by testcase
	Name string
	// Testcase function implement
	Implement reflect.Value
}

type TestSuite struct {
	// Name of testsuit, for example, LinuxSpec.Linux.Namespaces
	Name       string
	TestCase   []*TestCase
	TestResult []string
}

const (
	PASSED     = "passed"
	FAILED     = "failed"
	UNSPPORTED = "unspported"
	UNKNOWNERR = "unknowErr"
)

// TestResult to conver to json output
type TestResult struct {
	TestCaseName string `json:"testcasename"`
	// Json string : input of config
	Input interface{} `json:"input"`
	//funtion return error
	Err string `json:"error,omitempty"`
	//test result,passed,failed or unspported
	Result string `json:"result"`
}

var Manager *TestManager = new(TestManager)

// Add testSuite to TestManger
func (this *TestManager) AddTestSuite(testSuite TestSuite) {
	for _, ts := range this.TestSuite {
		if ts.Name == testSuite.Name {
			log.Fatalf("Existing same testsuite : %v", ts.Name)
		}
	}
	this.TestSuite = append(this.TestSuite, &testSuite)
}

// Set the TestResult structure
// testCaseName : TestCase.Name
// input : obj of test config, for example, LinuxSpec.Linux.Namespaces
// err : return value of TestCase.Implement func
// result : TestSuite.TestResult
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

// Conver TestResult to json string
func (this *TestResult) Marshal() string {
	js, err := json.Marshal(*this)
	if err != nil {
		log.Fatalf("Marshal error,%v\n", err)
	}
	return string(js)
}

// Add testcase to TestSuite
// impliment : TestCase.Implement, function of testcase
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

// Run testcases in TestSuite
func (this *TestSuite) Run() {

	for _, tc := range this.TestCase {
		var rt []reflect.Value
		rt = tc.Implement.Call(nil)
		str := rt[0].String()
		if str != "" {
			fmt.Printf("%-50s    %20s\n", tc.Name, getResult(str))
			this.TestResult = append(this.TestResult, str)
		}
	}
}

func getResult(str string) string {

	splitStr := "\"result\":\""
	retStr := strings.SplitAfter(str, splitStr)

	result := strings.TrimSuffix(retStr[1], "\"}")
	return result
}

// Merge jsonstring of each testcases in TestSuite into one json string
func (this *TestSuite) GetResult() string {

	result := "{\"" + this.Name + "\":["
	for _, tr := range this.TestResult {
		result = result + tr + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "]}"
	return result
}

func (this *TestManager) GetTotalResult() string {
	result := "{\" linuxspec  \":["
	for _, ts := range this.TestSuite {
		result = result + ts.GetResult() + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "]}"
	return result
}
