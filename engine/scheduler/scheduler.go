package main

import (
	"../lib/libocit"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

type ServerConfig struct {
	TSurl string
	CPurl string
	Debug bool
}

//public variable
var pub_conf ServerConfig
var pub_casename string
var pub_debug bool
var pub_casedir string

func ts_validation(ts_demo libocit.TestCase) (validate bool, err_string string) {
	if len(ts_demo.Name) > 0 {
	} else {
		err_string = "Cannot find the name"
		return false, err_string
	}

	if len(ts_demo.Requires) > 0 {
	} else {
		err_string = "Cannot find the libocit.Required resource"
		return false, err_string
	}
	return true, "OK"
}

func get_url(req libocit.Require, path string) (apiurl string) {
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

	u, _ := url.ParseRequestURI(apiuri)
	u.Path = path
	u.RawQuery = data.Encode()
	apiurl = fmt.Sprintf("%v", u)

	return apiurl
}

func apply_os(req libocit.Require) (resource libocit.Resource) {
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
			fmt.Println("err in read os")
			resource.ID = ""
			resource.Msg = "err in read os"
			resource.Status = false
		} else {
			if pub_debug {
				fmt.Println("Get OS reply ", string(body))
			}
			json.Unmarshal([]byte(body), &resource)
			resource.Req = req
			fmt.Println(resource)
		}
	}

	return resource
}

func apply_container(req libocit.Require) (resource libocit.Resource) {
	tar_url := libocit.TarFilelist(req.Files, pub_casedir, req.Class)
	post_url := pub_conf.CPurl + "/upload"
	libocit.SendFile(post_url, tar_url, tar_url)

	apiurl := pub_conf.CPurl + "/build"
	b, jerr := json.Marshal(req)
	if jerr != nil {
		fmt.Println("Failed to marshal json:", jerr)
		return
	}
	libocit.SendCommand(apiurl, []byte(b))
	return resource
}

func setContainerClass(deploys []libocit.Deploy, req libocit.Require) {
	for index := 0; index < len(deploys); index++ {
		deploy := deploys[index]
		for c_index := 0; c_index < len(deploy.Containers); c_index++ {
			if deploy.Containers[c_index].Class == req.Class {
				deploy.Containers[c_index].Distribution = req.Distribution
				deploy.Containers[c_index].Version = req.Version
			}
		}
	}
}

func apply_resources(ts_demo libocit.TestCase) (resources []libocit.Resource) {
	for index := 0; index < len(ts_demo.Requires); index++ {
		var resource libocit.Resource
		req := ts_demo.Requires[index]
		if req.Type == "os" {
			resource = apply_os(req)
		} else if req.Type == "container" {
			//FIXME: change the democase
			resource = apply_container(req)
			setContainerClass(ts_demo.Deploys, req)
		} else {
			fmt.Println("What is the new type? How can it pass the validation test")
		}
		resource.Used = false
		if len(resource.ID) > 1 {
			resources = append(resources, resource)
		}
	}
	return resources
}

func ar_validation(ar_demo []libocit.Resource) (validate bool, err_string string) {
	return true, "OK"
}

func main() {
	var ts_demo libocit.TestCase
	var validate bool
	var msg string
	var case_file string

	config_content := libocit.ReadFile("./scheduler.conf")
	json.Unmarshal([]byte(config_content), &pub_conf)

	pub_debug = pub_conf.Debug
	arg_num := len(os.Args)
	if arg_num < 2 {
		case_file = "./case01/Network-iperf.json"
	} else {
		case_file = os.Args[1]
	}
	pub_casedir = path.Dir(case_file)
	fmt.Println(case_file)
	test_json_str := libocit.ReadFile(case_file)
	json.Unmarshal([]byte(test_json_str), &ts_demo)
	if pub_debug {
		fmt.Println(ts_demo)
	}
	validate, msg = ts_validation(ts_demo)
	if !validate {
		fmt.Println(msg)
		return
	}
	if pub_debug {
		fmt.Println(ts_demo)
	}

	//libocit.Require Session
	var resources []libocit.Resource
	//TODO: async in the future
	resources = apply_resources(ts_demo)
	validate, msg = ar_validation(resources)
	if !validate {
		fmt.Println(msg)
		return
	}

	//Deploy Session

	// Prepare deploys
	for index := 0; index < len(ts_demo.Deploys); index++ {
		var deploy libocit.Deploy
		deploy = ts_demo.Deploys[index]
		for r_index := 0; r_index < len(resources); r_index++ {
			var resource libocit.Resource
			resource = resources[r_index]
			if resource.Used {
				continue
			}
			if resource.Req.Class == deploy.Class {
				ts_demo.Deploys[index].ResourceID = resource.ID
				resources[r_index].Used = true
				continue
			}
			// TODO should do it after apply resource
			fmt.Println("Cannot get here, failed to get enough resource")
		}
	}
	if pub_debug {
		fmt.Println(ts_demo.Deploys)
	}

	// Send deploys
	for index := 0; index < len(ts_demo.Deploys); index++ {
		if len(ts_demo.Deploys[index].ResourceID) > 0 {
			filelist := GetDeployFiles(ts_demo.Deploys[index])
			//FIXME: change the democase
			tar_url := libocit.TarFilelist(filelist, pub_casedir, ts_demo.Deploys[index].Object)
			post_url := pub_conf.TSurl + "/casefile/" + ts_demo.Deploys[index].ResourceID
			fmt.Println("Send file  -- ", post_url, tar_url)
			libocit.SendFile(post_url, tar_url, tar_url)

			apiurl := pub_conf.TSurl + "/deploy"
			b, jerr := json.Marshal(ts_demo.Deploys[index])
			if jerr != nil {
				fmt.Println("Failed to marshal json:", jerr)
				return
			}
			fmt.Println("Send command ", apiurl)
			libocit.SendCommand(apiurl, []byte(b))
		}
	}

	// Send 'Run' -- do we really need this?

	// Prepare collect
	for index := 0; index < len(ts_demo.Collects); index++ {
		for r_index := 0; r_index < len(ts_demo.Deploys); r_index++ {
			if ts_demo.Collects[index].Object == ts_demo.Deploys[r_index].Object {
				ts_demo.Collects[index].ResourceID = ts_demo.Deploys[r_index].ResourceID
			}
		}
	}

	// Send collects
	for index := 0; index < len(ts_demo.Collects); index++ {
		if len(ts_demo.Collects[index].ResourceID) > 0 {
			collect := ts_demo.Collects[index]
			for f_index := 0; f_index < len(collect.Files); f_index++ {
				file := collect.Files[f_index]
				apiurl := pub_conf.TSurl + "/report/" + ts_demo.Collects[index].ResourceID + "?file=" + file
				fmt.Println("Send collect cmd ", apiurl)
				resp, err := http.Get(apiurl)
				if err != nil {
					fmt.Println("Error ", err)
					continue
				}
				defer resp.Body.Close()
				resp_body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return
				}
				fmt.Println(resp.Status)
				fmt.Println(string(resp_body))
				cache_dir := "/tmp/test_scheduler_result"
				real_url := libocit.PreparePath(cache_dir, file)
				f, err := os.Create(real_url)
				defer f.Close()
				f.Write(resp_body)
				f.Sync()

			}
		}
	}

}

func GetDeployFiles(deploy libocit.Deploy) (filelist []string) {
	for index := 0; index < len(deploy.Files); index++ {
		filelist = append(filelist, deploy.Files[index])
	}
	for index := 0; index < len(deploy.Containers); index++ {
		container := deploy.Containers[index]
		for c_index := 0; c_index < len(container.Files); c_index++ {
			filelist = append(filelist, container.Files[c_index])
		}
	}
	return filelist
}
