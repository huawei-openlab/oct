package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
	"io/ioutil"
	"bytes"
	"os"
	"strconv"
)

type Require struct {
	Class string
	Type string
	Distribution string
	Version int
}

type Container struct {
	Object string
	Class string
	Cmd string
}

type Deploy struct {
	Object string
	Class string
	Cmd string
	Containers []Container

	ResourceID string
}

type Resource struct {
//TODO: put following to a struct and make a hash?
	ID  string	//returned 
	Status bool	//whether it is available
	Msg string	//return value from server

	Req Require
	_used bool
}

type TestCase struct {
	Name string
	License string
	Group string
	Requires []Require
	Deploys []Deploy
}

type ServerConfig struct {
	TSurl string
	CPurl string
}
//public variable
var pub_conf ServerConfig

var pub_debug bool

func parse(ts_str string) (ts_demo TestCase) {
	json.Unmarshal([]byte(ts_str), &ts_demo)

	return ts_demo
}

func ts_validation(ts_demo TestCase) (validate bool, err_string string){
	if len(ts_demo.Name) > 0 {
	} else {
		err_string = "Cannot find the name"
		return false, err_string
	}

	if len(ts_demo.Requires) > 0 {
	} else {
		err_string = "Cannot find the Required resource"
		return false, err_string
	}
	return true, "OK"
}

func debug_ts(ts_demo TestCase) {
	if !pub_debug {
		return
	}
	fmt.Println(ts_demo.Name)
	fmt.Println(ts_demo.Group)
	fmt.Println(ts_demo.Requires)
}

func read_conf()(config ServerConfig) {
	config_file := "./scheduler.conf"
	file, err := os.Open(config_file)
	defer file.Close()
	if err != nil {
	        fmt.Println(config_file, err)
	        return
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	json.Unmarshal([]byte(buf.String()), &config)
//	fmt.Println(config.TSurl, " ", config.CPurl)

	return config
}

func get_url(req Require, path string) (apiurl string) {
	var apiuri string
	data := url.Values{}
	if req.Type == "os" {
		apiuri = pub_conf.TSurl
	} else {
		apiuri = pub_conf.CPurl
	}
	if len(req.Distribution) > 1 {
		data.Add("Distribution", req.Distribution)
	}
	data.Add("Version", strconv.Itoa(req.Version))

        u, _:= url.ParseRequestURI(apiuri)
        u.Path = path
        u.RawQuery = data.Encode()
        apiurl = fmt.Sprintf("%v", u)

	return apiurl
}

func apply_os(req Require) (resource Resource){
	var apiurl string

	apiurl = get_url(req, "/os")
	if pub_debug {
		fmt.Println("get url: ", apiurl)
	}
	resp, err := http.Get(apiurl)
	defer resp.Body.Close()
	if err != nil {
		// handle error
		fmt.Println("err in get")
		resource.ID = ""
		resource.Msg = "err in get os"
		resource.Status = false
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resource.ID = ""
			resource.Msg = "err in read os"
			resource.Status = false
		} else {
			json.Unmarshal([]byte(body), &resource)
			resource.Req = req
			fmt.Println(resource)
		}	
	}

	return resource
}

func apply_container(req Require) (resource Resource){
	return resource
}

func apply_resource(req Require) (resource Resource){
	if req.Type == "os" {
		resource = apply_os(req)
	} else if req.Type == "container" {
		resource = apply_container(req)
	} else {
		fmt.Println("What is the new type? How can it pass the validation test")
	}
	resource._used = false
	return resource
}

func apply_resources(ts_demo TestCase) (resources []Resource){
	for index :=0; index < len(ts_demo.Requires); index++ {
		var resource Resource
		var req Require
		req = ts_demo.Requires[index]
		resource = apply_resource(req)
		if len(resource.ID) > 1 {
			resources = append(resources, resource)	
		}
	}
	return resources
}

func ar_validation(ar_demo []Resource) (validate bool, err_string string){
	return true, "OK"
}

func debug_ar(ar_demo []Resource) {
	if !pub_debug {
		return
	}

	fmt.Println(ar_demo)
}

func read_testcase(ts_file string) (testcase string) {
	file, err := os.Open(ts_file)
	defer file.Close()
	if err != nil {
	        fmt.Println(ts_file, err)
	        return
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	testcase = buf.String()
	
	if pub_debug {
		fmt.Println(testcase)
	}

	return testcase
}

func debug_deploy(deploys []Deploy) {
	if pub_debug {
		fmt.Println("Debug deploys ", deploys)
	}
}

//send deploy to the testserver and the server will do the deploy work
func send_deploy(deploy Deploy) {
	var apiurl string

	apiurl = pub_conf.TSurl + "/deploy"
	if pub_debug {
		fmt.Println("get url: ", apiurl)
	}
	b, jerr := json.Marshal(deploy)
	if jerr != nil {
		fmt.Println("Failed to marshal json:", jerr)
		return
	}
fmt.Println(deploy)
	body := bytes.NewBuffer([]byte(b))
	resp, perr := http.Post(apiurl, "application/json;charset=utf-8", body)
	defer resp.Body.Close()
	if perr != nil {
		// handle error
		fmt.Println("err in post:", perr)
		return
	} else {
		result, berr := ioutil.ReadAll(resp.Body)
		if berr != nil {
		} else {
//TODO: finish the result check
			if pub_debug {
				fmt.Println(result)
			}
		}	
	}
}

func main() {
	var ts_demo TestCase
	var validate bool
	var msg string
	var test_json_str string

	pub_debug = false

	pub_conf = read_conf()
	test_json_str = read_testcase("./democase.json")
	ts_demo = parse(test_json_str)
	if pub_debug {
		fmt.Println(ts_demo)
	}
	validate, msg = ts_validation(ts_demo)
	if !validate {
		fmt.Println(msg)
		return
	}
	debug_ts(ts_demo)

//Require Session
	var resources []Resource
//TODO: async in the future
	resources = apply_resources(ts_demo)
	validate, msg = ar_validation(resources)
	if !validate {
		fmt.Println(msg)
		return
	}
	debug_ar(resources)

//Deploy Session

// Prepare deploys
	for index :=0; index < len(ts_demo.Deploys); index++ {
		var deploy Deploy
		deploy = ts_demo.Deploys[index]
		for r_index := 0; r_index < len(resources); r_index++ {
			var resource Resource
			resource = resources[r_index]
			if resource._used {
				continue
			}
			if resource.Req.Class == deploy.Class {
				ts_demo.Deploys[index].ResourceID = resource.ID
				resources[r_index]._used = true
				continue
			}
// TODO should do it after apply resource
			fmt.Println("Cannot get here, failed to get enough resource")
		}
	}
	debug_deploy(ts_demo.Deploys)

// Send deploys
	for index :=0; index < len(ts_demo.Deploys); index++ {
		if len(ts_demo.Deploys[index].ResourceID) > 0 {
			send_deploy(ts_demo.Deploys[index])
		}
	}

}

