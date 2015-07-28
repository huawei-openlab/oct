package main

import (
	"fmt"
	"io"
	"os"
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

func uploadFile(w http.ResponseWriter, r *http.Request){
	if "POST" == r.Method {
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
		return
	}
	
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

func main() {
	var config OCITDConfig
	config = read_conf()
	var port string
	port = fmt.Sprintf(":%d", config.Port)

	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/os", testOS)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
