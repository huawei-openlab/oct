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
		err_string = "Cannot find the libocit.Requires resource"
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
	data.Add("Version", req.Version)

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
	var files []string
	for index := 0; index < len(req.Files); index++ {
		files = append(files, path.Join(pub_casedir, req.Files[index]))
	}
	tar_url := libocit.TarFilelist(files, pub_casedir, req.Class)
	post_url := pub_conf.CPurl + "/upload"

	var params map[string]string
	libocit.SendFile(post_url, tar_url, params)

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

//Usage:  ./scheduler ./demo.tar.gz
func main() {
	var case_file string

	config_content := libocit.ReadFile("./scheduler.conf")
	json.Unmarshal([]byte(config_content), &pub_conf)

	pub_debug = pub_conf.Debug
	arg_num := len(os.Args)
	if arg_num < 2 {
		case_file = "./demo.tar.gz"
	} else {
		case_file = os.Args[1]
	}
	post_url := pub_conf.TSurl + "/task"
	fmt.Println(post_url)

	var params map[string]string
	params = make(map[string]string)
	//FIXME: the id autobe automaticly allocated
	params["id"] = "00001"
	libocit.SendFile(post_url, case_file, params)
	return

}

func GetDeployFiles(dir string, deploy libocit.Deploy) (filelist []string) {
	for index := 0; index < len(deploy.Files); index++ {
		filelist = append(filelist, path.Join(dir, deploy.Files[index]))
	}
	for index := 0; index < len(deploy.Containers); index++ {
		container := deploy.Containers[index]
		for c_index := 0; c_index < len(container.Files); c_index++ {
			filelist = append(filelist, path.Join(dir, container.Files[c_index]))
		}
	}
	return filelist
}
