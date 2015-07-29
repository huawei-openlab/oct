package main

import (
        "github.com/drones/routes"
	"fmt"
        "archive/tar"
        "compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"net/http"
	"log"
	"bytes"
	"encoding/json"
	"path"
)

func preparePath(filename string) (realurl string) {
//TODO: should add the 'testcase name'
	pre_uri := "/tmp/testcase_ocitd_cache/"
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

func testOS(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.URL.Query().Get("Distribution"))
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
	fmt.Println(handler.Filename)
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
	io.WriteString(w, "OK")

	untarFiles(real_url)
	return
	
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	html := `
	<form enctype="multipart/form-data" action="/upload" method="POST">
		Send this file : <input name="tsfile" type="file" /?
		<input type="submit" value="Send File"/>
	</form>
	`

	io.WriteString(w, html)
}

type OCITDConfig struct {
	TSurl string
	Port int
}
 
func read_conf() (config OCITDConfig) {
	config_file := "./ocitd.conf"
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

type Container struct {
        Object string
        Class string
        Cmd string
	Files []string
	Distribution string
	Version int
}

type Deploy struct {
        Object string
        Class string
        Cmd string
	Files []string
        Containers []Container

        ResourceID string
}

func RunCommand(cmd string) {
	pre_uri := "/tmp/testcase_ocitd_cache/"
	os.Chdir(pre_uri)
	c := exec.Command("/bin/sh", "-c", cmd)
	c.Run()
}

func PullImage(container Container) {
	if container.Distribution == "Docker" {
		cmd := "docker pull " + container.Class
		c := exec.Command("/bin/sh", "-c", cmd)
		c.Run()

		fmt.Println("Exec pull image ", cmd)
	}
}

func DeployCommand(w http.ResponseWriter, r *http.Request){
        result, _:= ioutil.ReadAll(r.Body)
        r.Body.Close()

        var deploy Deploy
        json.Unmarshal([]byte(result), &deploy)

	images := make(map[string] int)
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

func main() {
	var config OCITDConfig
	config = read_conf()
	var port string
	port = fmt.Sprintf(":%d", config.Port)

	mux := routes.New()
        mux.Post("/upload", UploadFile)
        mux.Post("/deploy", DeployCommand)
	mux.Get("/os", testOS)
        http.Handle("/", mux)
        err := http.ListenAndServe(port, nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}
