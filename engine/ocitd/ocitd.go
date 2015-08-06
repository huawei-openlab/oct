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

type OCITDConfig struct {
	TSurl string
	Port  int
}

func GetReport(w http.ResponseWriter, r *http.Request) {
	pre_uri := "/tmp/testcase_ocitd_cache"
	filename := r.URL.Query().Get("file")
	realurl := path.Join(pre_uri, filename)

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
	cache_uri := "/tmp/testcase_ocitd_cache"
	real_url, _ := libocit.ReceiveFile(w, r, cache_uri)

	libocit.UntarFile(cache_uri, real_url)
	w.Write([]byte("OK"))
	return
}

func RunCommand(cmd string) {
	pre_uri := "/tmp/testcase_ocitd_cache/source"
	os.Chdir(pre_uri)

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

func DeployCommand(w http.ResponseWriter, r *http.Request) {
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var deploy libocit.Deploy
	json.Unmarshal([]byte(result), &deploy)

	images := make(map[string]int)
	for index := 0; index < len(deploy.Containers); index++ {
		container := deploy.Containers[index]
		_, pulled := images[container.Class]
		if pulled {
			fmt.Println("Already pulled")
		} else {
			images[container.Class] = 1
			PullImage(container)
		}
	}

	if len(deploy.Cmd) > 0 {
		RunCommand(deploy.Cmd)
	}
}

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
	var config OCITDConfig
	content := libocit.ReadFile("./ocitd.conf")
	json.Unmarshal([]byte(content), &config)

	var port string
	port = fmt.Sprintf(":%d", config.Port)

	mux := routes.New()
	mux.Get("/report", GetReport)
	mux.Post("/upload", UploadFile)
	mux.Post("/deploy", DeployCommand)
	mux.Post("/collect", CollectFiles)
	http.Handle("/", mux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
