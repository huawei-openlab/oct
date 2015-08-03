package main

import (
	"../lib/libocit"
	"../lib/routes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type ContainerPoolConfig struct {
	Port  int
	Debug bool
}

func main() {
	var config ContainerPoolConfig

	content := libocit.ReadFile("./containerpool.conf")
	json.Unmarshal([]byte(content), &config)
	var port string
	port = fmt.Sprintf(":%d", config.Port)

	mux := routes.New()
	mux.Post("/upload", UploadFile)
	mux.Post("/build", BuildImage)
	http.Handle("/", mux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var lock = sync.RWMutex{}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("tcfile")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()

	cache_uri := "/tmp/testcase_ocitd_cache"
	real_url := libocit.PreparePath(cache_uri, handler.Filename)
	f, err := os.Create(real_url)
	if err != nil {
		//TODO: better system error
		http.Error(w, err.Error(), 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	libocit.UntarFile(cache_uri, real_url)

	w.Write([]byte("OK"))
}

func BuildImage(w http.ResponseWriter, r *http.Request) {
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var req libocit.Require
	json.Unmarshal([]byte(result), &req)

	fmt.Println(req)
}
