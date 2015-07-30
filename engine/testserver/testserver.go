package main

import (
	"../lib/routes"
	"log"
	"net/http"
	"sync"
	"strconv"
	"os"
	"io"
	"bytes"
	"fmt"
	"encoding/json"
	"mime/multipart"
	"io/ioutil"
	"path"
)

type TestServerConfig struct{
	Port int
	ServerListFile string
	Debug bool
}

func read_conf()(config TestServerConfig) {
	config_file := "./testserver.conf"
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
	pre_uri := "/tmp/testserver_cache/"
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

func send_deploy_file(id string, realurl string, filename string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("tsfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return
	}
	fh, err := os.Open(realurl)
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	os := *(store[id])

//FIXME: better org
	post_url := "http://" + os.IP + ":9001/upload"
	fmt.Println(post_url)
	resp, err := http.Post(post_url, contentType, bodyBuf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
        fmt.Println(resp.Status)
        fmt.Println(string(resp_body))
}

//send deploy to the ocitd and the client will do the deploy work
func SendDeployCommand(deploy Deploy) {
        var apiurl string
	os := *(store[deploy.ResourceID])
	deploy.ResourceID = ""
	apiurl = "http://" + os.IP + ":9001/deploy"
        b, jerr := json.Marshal(deploy)
        if jerr != nil {
                fmt.Println("Failed to marshal json:", jerr)
                return
        }

        body := bytes.NewBuffer([]byte(b))
        resp, perr := http.Post(apiurl, "application/json;charset=utf-8", body)
        defer resp.Body.Close()
        if perr != nil {
                // handle error
                fmt.Println("err in post:", perr)
                return
        } else {
                result, berr := ioutil.ReadAll(resp.Body)
                if berr != nil {
                } else {
                                fmt.Println(result)
                }
        }
}

func DeployCommand(w http.ResponseWriter, r *http.Request){
	result, _:= ioutil.ReadAll(r.Body)
	r.Body.Close()

	var deploy Deploy
	json.Unmarshal([]byte(result), &deploy)
	
	SendDeployCommand(deploy)
//FIXME: write back to the scheduler
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
	id := r.URL.Query().Get(":ID")

//TS act as the agent
	send_deploy_file(id, real_url, handler.Filename)
	w.Write([]byte(id))
}

func main() {
	var config TestServerConfig

	config = read_conf()
	var port string
	port = fmt.Sprintf(":%d", config.Port)
//	pub_debug = config.Debug
	init_db (config.ServerListFile)

	mux := routes.New()
	mux.Get("/os", GetOS)
	mux.Post("/os", PostOS)
	mux.Post("/upload/:ID", UploadFile)
	mux.Post("/deploy", DeployCommand)
	http.Handle("/", mux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type OS struct {
	ID string
	Distribution string
	Version string
	Arch string
	CPU int64
	Memory int64
	IP string
	locked bool
}

var store = map[string]*OS{}

var lock = sync.RWMutex{}

func GetOSQuery(r *http.Request) (os OS) {
	os.Distribution = r.URL.Query().Get("Distribution")
	os.Version = r.URL.Query().Get("Version")
	os.Arch = r.URL.Query().Get("Arch")

	var cpu string
	cpu = r.URL.Query().Get("CPU")
	if len(cpu) > 0 {
		cpu_count, cpu_err := strconv.ParseInt(cpu, 10, 64)
 		if cpu_err != nil {
			//TODO, should report the err
		} else {
			os.CPU = cpu_count
		}
	} else {
		os.CPU = 0
	}

	var memory string
	memory = r.URL.Query().Get("Memory")
	if len(memory) > 0 {
		memory_count, memory_err := strconv.ParseInt(cpu, 10, 64)
 		if memory_err != nil {
			//TODO, should report the err
		} else {
			os.Memory = memory_count
		}
	} else {
		os.Memory = 0
	}

	log.Println(os)
	return os
}

// Will use sql to seach, for now, just
func GetAvaliableResource(os_query OS) (ID string) {
	for _, os := range store {
		if len(os_query.Distribution) > 1 {
			if os_query.Distribution != (*os).Distribution {
				continue
			}
		}
		if len(os_query.Version) > 1 {
			if os_query.Version != (*os).Version {
				continue
			}
		}
		if len(os_query.Arch) > 1 {
			if os_query.Arch != (*os).Arch {
				continue
			}
		}
		if os_query.CPU >  (*os).CPU {
			log.Println("not enough CPU")
			continue
		}
		if os_query.Memory > (*os).Memory {
			log.Println("not enough Memory")
			continue
		}
		ID = (*os).ID
		return ID
	}
	return ""
}

func GetOS(w http.ResponseWriter, r *http.Request){
	var os_query OS
	os_query = GetOSQuery (r)
	if len(os_query.Distribution) < 1 {
		GetAllOS(w, r)
		return
	}
	lock.RLock()

	var ID string
	ID = GetAvaliableResource(os_query)
	lock.RUnlock()

	log.Println(ID)
	if len(ID) < 1 {
		return
	}

//FIXME, the struct like Resource should be in the lib
	type Resource struct {
		ID string
		Msg string
		Status bool
	}
	var resource Resource
	resource.ID = ID
	resource.Msg = "ok, good resource"
	resource.Status = true
	body, _:= json.Marshal(resource)
	w.Write([]byte(body))
}

func GetAllOS(w http.ResponseWriter, r *http.Request){
	lock.RLock()
	os_list := make([]OS, len(store))
	i := 0
	for _, os := range store {
		os_list[i] = *os
		i++
	}
	lock.RUnlock()
	w.Write([]byte("FIXME: all the os"))
}

func PostOS(w http.ResponseWriter, r *http.Request){
	var os OS
	result, _:= ioutil.ReadAll(r.Body)
	r.Body.Close()

	json.Unmarshal([]byte(result), &os)
	if os.Distribution == "" {
		w.Write([]byte("os distribution required"))
		return
	}
	if os.Version == "" {
		w.Write([]byte("os version required"))
		return
	}
	if os.Arch == "" {
		w.Write([]byte("os arch required"))
		return
	}
	lock.Lock()
	store[os.Distribution] = &os
	lock.Unlock()
}

func DeleteOS(w http.ResponseWriter, r *http.Request){
	Distribution := r.URL.Query().Get("Distribution")
	lock.Lock()
	delete(store, Distribution)
	lock.Unlock()
	w.WriteHeader(http.StatusOK)
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

func sendCommand(ID string, CMD string) {
	fmt.Println(ID, CMD)
}

func deployRequest(deploy Deploy) {
	os := *(store[deploy.ResourceID])
	fmt.Println("the deploy request is: ", os)
	if len(deploy.Cmd) > 0 {
		sendCommand(deploy.ResourceID, deploy.Cmd)
	}
}

// Will use DB in the future, (mongodb for example)
func init_db (filename string) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		fmt.Println(filename, err)
		return
	}

	type _Servers struct{
		Servers []OS
	}
	var _servers _Servers
	
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	json.Unmarshal([]byte(buf.String()), &_servers)

	for index := 0; index < len(_servers.Servers); index++ {
		os := _servers.Servers[index]
		os.locked = false
		store[os.ID] = &os
	}

	fmt.Println(store)
}
