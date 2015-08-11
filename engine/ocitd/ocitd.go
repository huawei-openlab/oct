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
	filename := r.URL.Query().Get("File")
	ID := r.URL.Query().Get("ID")
	realurl := path.Join(pub_config.CacheDir, ID, filename)

	file, err := os.Open(realurl)
	defer file.Close()
	if err != nil {
		//FIXME: add to head
		w.Write([]byte("Cannot open the file: " + realurl))
		return
	}

	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)

	fmt.Println(realurl)

	w.Write([]byte(buf.String()))

}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	//The task ID is alreay included in the real_url
	real_url, _ := libocit.ReceiveFile(w, r, pub_config.CacheDir)

	libocit.UntarFile(pub_config.CacheDir, real_url)

	var ret libocit.HttpRet
	ret.Status = "OK"
	ret_string, _ := json.Marshal(ret)
	w.Write([]byte(ret_string))

	return
}

func RunCommand(cmd string) {
	os.Chdir(path.Join(pub_config.CacheDir, "source"))

	C.CSystem(C.CString(cmd))
	return
	fmt.Println("Run the command ", cmd)
	c := exec.Command("/bin/sh", "-c", cmd)
	c.Run()
	fmt.Println("After run the command ", cmd)
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

//FIXME: used in OCTD, it is necessary to put here?
func FindCommand(testCase libocit.TestCase, objectName string, phaseName string) (command string) {
	var deploy_list []libocit.Deploy
	if phaseName == "deploy" {
		deploy_list = testCase.Deploys
	} else if phaseName == "run" {
		deploy_list = testCase.Run
	}
	for index := 0; index < len(deploy_list); index++ {
		if deploy_list[index].Object == objectName {
			command = deploy_list[index].Cmd
			break
		}
	}
	return command
}

func UpdateStatus(testCommand libocit.TestingCommand) {
	var testStatus libocit.TestingStatus

	post_url := pub_config.TSurl + "/" + testCommand.ID + "/status"
	if testCommand.Command == "deploy" {
		testStatus.Status = "Deployed"
	} else if testCommand.Command == "run" {
		testStatus.Status = "Finish"
	}
	testStatus.Object = testCommand.Object
	ts_string, _ := json.Marshal(testStatus)
	libocit.SendCommand(post_url, []byte(ts_string))
}

func TestingCommand(w http.ResponseWriter, r *http.Request) {
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var testCommand libocit.TestingCommand
	json.Unmarshal([]byte(result), &testCommand)

	var testCase libocit.TestCase
	content := libocit.ReadFile(path.Join(pub_config.CacheDir, testCommand.ID, "config.json"))
	json.Unmarshal([]byte(content), &testCase)

	command := FindCommand(testCase, testCommand.Object, testCommand.Command)

	if len(command) > 0 {
		RunCommand(command)
	}
	//Send status update to the test server
	UpdateStatus(testCommand)

	var ret libocit.HttpRet
	ret.Status = "OK"
	ret_string, _ := json.Marshal(ret)
	w.Write([]byte(ret_string))
}

//FIXME: should be removed, since move to test server
func CollectFiles(w http.ResponseWriter, r *http.Request) {
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var filelist []string
	json.Unmarshal([]byte(result), &filelist)

	for index := 0; index < len(filelist); index++ {
		fmt.Println(filelist[index])
	}

}

func main() {
	content := libocit.ReadFile("./ocitd.conf")
	json.Unmarshal([]byte(content), &pub_config)

	var port string
	port = fmt.Sprintf(":%d", pub_config.Port)

	mux := routes.New()
	mux.Get("/result", GetResult)
	mux.Post("/task", UploadFile)
	mux.Post("/command", TestingCommand)
	//	mux.Post("/collect", CollectFiles)
	http.Handle("/", mux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
