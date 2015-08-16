package main

import (
	"../lib/libocit"
	"encoding/json"
	"fmt"
	"os"
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
