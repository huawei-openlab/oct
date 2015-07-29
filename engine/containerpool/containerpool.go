package main

import (
	"github.com/drones/routes"
	"log"
	"net/http"
	"sync"
	"os"
	"archive/tar"
	"compress/gzip"
	"io"
	"bytes"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path"
)

type ContainerPoolConfig struct{
	Port int
	Debug bool
}

func read_conf()(config ContainerPoolConfig) {
	config_file := "./containerpool.conf"
	file, err := os.Open(config_file)
	defer file.Close()
	if err != nil {
		fmt.Println(config_file, err)
		return
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	json.Unmarshal([]byte(buf.String()), &config)

	return config
}

func preparePath(filename string) (realurl string) {
//TODO: should add the 'testcase name'
	pre_uri := "/tmp/containerpool_cache/"
	realurl = path.Join(pre_uri, filename)
	dir := path.Dir(realurl)
	p, err:= os.Stat(dir)
	if err!= nil {
		if !os.IsExist(err) {
			os.MkdirAll(dir, 0777)
		}
	} else {
		if p.IsDir() {
			return realurl
		} else {
			os.Remove(dir)
			os.MkdirAll(dir, 0777)
		}
	}
	return realurl
}


func main() {
	var config ContainerPoolConfig

	config = read_conf()
	var port string
	port = fmt.Sprintf(":%d", config.Port)
//	pub_debug = config.Debug

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

type Require struct {
	Class string
	Type string
	Distribution string
	Version int
	Files []string
}


func untarFiles(filename string) {
	fr, err:= os.Open(filename)
	if err != nil {
		fmt.Println("fail in open file ", filename)
		return
	}
	defer fr.Close()
	gr, err:= gzip.NewReader(fr)
	if err != nil {
		fmt.Println("fail in using gzip")
		return
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic (err)
		}
		real_name := preparePath(h.Name)
		fw, _:= os.OpenFile(real_name, os.O_CREATE | os.O_WRONLY, os.FileMode(h.Mode))
		defer fw.Close()

		io.Copy(fw, tr)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("tsfile")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	real_url := preparePath(handler.Filename)
	f,err:=os.Create(real_url)
	if err != nil {
//TODO: better system error
		http.Error(w, err.Error(), 500)
		return
	}
	defer f.Close()
	io.Copy(f,file)
	untarFiles(real_url)

	w.Write([]byte("OK"))
}

func BuildImage(w http.ResponseWriter, r *http.Request){
	result, _:= ioutil.ReadAll(r.Body)
	r.Body.Close()

	var req Require
	json.Unmarshal([]byte(result), &req)

	fmt.Println(req)
}

type Container struct {
	Object string
	Class string
	Cmd string
	Files []string
	Distribution string
	Version int
}

//TODO add a 'casename' ?
type Deploy struct {
	Object string
	Class string
	Cmd string
	Files []string
	Containers []Container

	ResourceID string
}

