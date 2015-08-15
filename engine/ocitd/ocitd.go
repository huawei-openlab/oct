package main

import (
	"../lib/libocit"
	"../lib/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

//	"strings"
)

/*
#include <stdio.h>
#include <stdlib.h>
int CSystem(char *cmd){
	return system (cmd);
}
*/
import "C"

type OCTDConfig struct {
	TSurl    string
	Port     int
	CacheDir string
	Debug    bool
}

var pub_config OCTDConfig

func GetResult(w http.ResponseWriter, r *http.Request) {
	var realurl string
	filename := r.URL.Query().Get("File")
	ID := r.URL.Query().Get("ID")

	_, err := os.Stat(filename)
	if err == nil {
		//absolute path
		realurl = filename
	} else {
		json_dir := FindJsonDir(path.Join(pub_config.CacheDir, ID))
		realurl = path.Join(json_dir, filename)
		_, err = os.Stat(realurl)
		if err != nil {
			w.Write([]byte("Cannot find the file: " + realurl))
			return
		}
	}
	file, err := os.Open(realurl)
	defer file.Close()
	if err != nil {
		//FIXME: add to head
		w.Write([]byte("Cannot open the file: " + realurl))
		return
	}

	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)

	w.Write([]byte(buf.String()))
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	real_url, params := libocit.ReceiveFile(w, r, pub_config.CacheDir)

	fmt.Println(params)

	if val, ok := params["id"]; ok {
		libocit.UntarFile(path.Join(pub_config.CacheDir, val), real_url)
	} else {

		libocit.UntarFile(pub_config.CacheDir, real_url)
	}
	var ret libocit.HttpRet
	ret.Status = "OK"
	ret_string, _ := json.Marshal(ret)
	w.Write([]byte(ret_string))

	return
}

func RunCommand(cmd string, dir string) {
	if pub_config.Debug {
		fmt.Println("Run the command < ", cmd, ">  in ", dir)
	}
	os.Chdir(dir)

	debugging := true
	if debugging {
		c := exec.Command("/bin/sh", "-c", cmd)
		c.Run()
	} else {
		C.CSystem(C.CString(cmd))
	}
	return

	// Golang bug? cannot get the standard output
	//	fmt.Println("Run the command ", cmd)
	//	c := exec.Command("/bin/sh", "-c", cmd)
	//	c.Run()
	//	fmt.Println("After run the command ", cmd)
}

func PullImage(container libocit.Container) {
	//FIXME: no need to do this!
	return
	if container.Distribution == "Docker" {
		cmd := "docker pull " + container.Class
		c := exec.Command("/bin/sh", "-c", cmd)
		c.Run()

		fmt.Println("Exec pull image ", cmd)
	}
}

func UpdateStatus(testCommand libocit.TestingCommand) {
	var testStatus libocit.TestingStatus

	post_url := pub_config.TSurl + "/" + testCommand.ID + "/status"
	if testCommand.Status == "deploy" {
		testStatus.Status = "Deployed"
	} else if testCommand.Status == "run" {
		testStatus.Status = "Finish"
	}
	testStatus.Object = testCommand.Object
	ts_string, _ := json.Marshal(testStatus)
	libocit.SendCommand(post_url, []byte(ts_string))
}

//This is for the un-formal testcase, for example with third-party libs included
//TODO, need to use the formal format, since the output will also be '.json'
func FindJsonDir(base_dir string) (json_dir string) {
	files_info, _ := ioutil.ReadDir(base_dir)
	for _, file := range files_info {
		if file.IsDir() {
			sub_json_dir := FindJsonDir(path.Join(base_dir, file.Name()))
			if len(sub_json_dir) > 1 {
				json_dir = sub_json_dir
				return json_dir
			}
		} else {
			fileSuffix := path.Ext(file.Name())
			if fileSuffix == ".json" {
				_, err := os.Stat(path.Join(base_dir, "source"))
				if err != nil {
					continue
				} else {
					//  ./casename.json, ./source/
					json_dir = base_dir
					return json_dir
				}
			}
		}
	}
	return json_dir
}

func TestingCommand(w http.ResponseWriter, r *http.Request) {
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var testCommand libocit.TestingCommand
	json.Unmarshal([]byte(result), &testCommand)

	if len(testCommand.Command) > 0 {
		json_dir := FindJsonDir(path.Join(pub_config.CacheDir, testCommand.ID))
		dir := path.Join(json_dir, "source")
		RunCommand(testCommand.Command, dir)
	}
	//Send status update to the test server
	UpdateStatus(testCommand)

	var ret libocit.HttpRet
	ret.Status = "OK"
	ret_string, _ := json.Marshal(ret)
	w.Write([]byte(ret_string))
}

func RegisterToTestServer() {
	post_url := pub_config.TSurl + "/os"

	//TODO
	//Seems there will be lots of coding while getting the system info
	//Using config now.

	content := libocit.ReadFile("./host.conf")
	fmt.Println(content)
	ret := libocit.SendCommand(post_url, []byte(content))
	fmt.Println(ret)
}

func main() {
	content := libocit.ReadFile("./ocitd.conf")
	json.Unmarshal([]byte(content), &pub_config)

	RegisterToTestServer()

	var port string
	port = fmt.Sprintf(":%d", pub_config.Port)

	mux := routes.New()
	mux.Get("/result", GetResult)
	mux.Post("/task", UploadFile)
	mux.Post("/command", TestingCommand)

	http.Handle("/", mux)
	fmt.Println("Start to listen ", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
