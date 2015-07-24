package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

type Require struct {
	Rtype string
	Rname string
	Version int
}

type Resource struct {
	Req Require
//TODO: put following to a struct and make a hash?
	ID  string	//returned 
	Status bool	//whether it is available
	Msg string	//return value from server
}

type TestCase struct {
	Name string
	License string
	Group string
	Requires []Require
}

const test_json_str = `{"Name": "Network-performance", "License": "GPL",
"Requires": [{"Rtype": "os", "Rname": "CentOS", "Version": 7}, {"Rtype": "container", "Rname": "Docker", "Version": 1}
]}`

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
	fmt.Println(ts_demo.Name)
	fmt.Println(ts_demo.Group)
	fmt.Println(ts_demo.Requires)
}

func apply_os(req Require) (resource Resource){
	var default_url string

	default_url = "http://localhost:8080/os?Distribution=Ubuntu"
	resp, err := http.Get(default_url)
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
			fmt.Println("return the os id from server " + string(body))
//			json.Unmarshal([]byte(body), &resource)
			resource.Req = req
			resource.ID = string(body)
			resource.Msg = "ok"
			resource.Status = true
			fmt.Println(resource)
		}	
	}

	return resource
}

func apply_container(req Require) (resource Resource){
	return resource
}

func apply_resource(req Require) (resource Resource){
	if req.Rtype == "os" {
		resource = apply_os(req)
	} else if req.Rtype == "container" {
		resource = apply_container(req)
	} else {
		fmt.Println("What is the new type? How can it pass the validation test")
	}
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
	fmt.Println(ar_demo)
}

func main() {
	var ts_demo TestCase
	var validate bool
	var msg string
        ts_demo = parse(test_json_str)
	validate, msg = ts_validation(ts_demo)
	if !validate {
		fmt.Println(msg)
		return
	}
	debug_ts(ts_demo)

	var resources []Resource
//TODO: async in the future
	resources = apply_resources(ts_demo)
	validate, msg = ar_validation(resources)
	if !validate {
		fmt.Println(msg)
		return
	}
	debug_ar(resources)
}

