package main

import (
	"../lib/libocit"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type ServerConfig struct {
	TSurl string
	CPurl string
	Debug bool
}

//public variable
var pub_config ServerConfig
var pub_casedir string

//TODO the following container function should move the the container service 
func apply_container(req libocit.Require) {
	var files []string
	for index := 0; index < len(req.Files); index++ {
		files = append(files, path.Join(pub_casedir, req.Files[index]))
	}
	tar_url := libocit.TarFilelist(files, pub_casedir, req.Class)
	post_url := pub_config.CPurl + "/upload"

	var params map[string]string
	libocit.SendFile(post_url, tar_url, params)

	apiurl := pub_config.CPurl + "/build"
	b, jerr := json.Marshal(req)
	if jerr != nil {
		fmt.Println("Failed to marshal json:", jerr)
		return
	}
	libocit.SendCommand(apiurl, []byte(b))
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

//Usage:  ./scheduler ./demo.tar.gz
func main() {
	var case_file string

	config_content := libocit.ReadFile("./scheduler.conf")
	json.Unmarshal([]byte(config_content), &pub_config)

	arg_num := len(os.Args)
	if arg_num < 2 {
		case_file = "./democase.tar.gz"
	} else {
		case_file = os.Args[1]
	}
	post_url := pub_config.TSurl + "/task"
	fmt.Println(post_url)

	var params map[string]string
	params = make(map[string]string)
	//TODO: use system time as the id now
	id := fmt.Sprintf("%d", time.Now().Unix())
	params["id"] = id
	ret := libocit.SendFile(post_url, case_file, params)
	fmt.Println(ret)
	return
}
