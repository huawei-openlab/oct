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
	"strconv"
	"sync"
)

type TestServerConfig struct {
	Port           int
	ServerListFile string
	Debug          bool
}

func DeployCommand(w http.ResponseWriter, r *http.Request) {
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var deploy libocit.Deploy
	json.Unmarshal([]byte(result), &deploy)

	os := *(store[deploy.ResourceID])
	deploy.ResourceID = ""
	deploy_url := "http://" + os.IP + ":9001/deploy"

	fmt.Println("Receive and send the deploy command ", deploy_url)
	libocit.SendCommand(deploy_url, []byte(result))
	//TODO: write back the info
}

func GetReport(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":ID")
	report_file := r.URL.Query().Get("file")
	os := *(store[id])

	//FIXME: better org
	get_url := "http://" + os.IP + ":9001/report?file=" + report_file
	log.Println("In get report: ", get_url)
	resp, err := http.Get(get_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Println(resp.Status)
	log.Println(string(resp_body))
	w.Write([]byte(string(resp_body)))
}

func ReceiveFile(w http.ResponseWriter, r *http.Request) {
	cache_uri := "/tmp/testserver_cache/"
	real_url, handle_name := libocit.ReceiveFile(w, r, cache_uri)

	id := r.URL.Query().Get(":ID")

	//TS act as the agent
	os := *(store[id])

	//FIXME: better org
	post_url := "http://" + os.IP + ":9001/upload"
	fmt.Println("Receive and send file ", real_url)
	libocit.SendFile(post_url, real_url, handle_name)
	w.Write([]byte(id))
}

func GetOSQuery(r *http.Request) (os libocit.OS) {
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
func GetAvaliableResource(os_query libocit.OS) (ID string) {
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
		if os_query.CPU > (*os).CPU {
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

func GetOS(w http.ResponseWriter, r *http.Request) {
	var os_query libocit.OS
	os_query = GetOSQuery(r)
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

	var resource libocit.Resource
	resource.ID = ID
	resource.Msg = "ok, good resource"
	resource.Status = true
	body, _ := json.Marshal(resource)
	w.Write([]byte(body))
}

func GetAllOS(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	os_list := make([]libocit.OS, len(store))
	i := 0
	for _, os := range store {
		os_list[i] = *os
		i++
	}
	lock.RUnlock()
	w.Write([]byte("FIXME: all the os"))
}

func PostOS(w http.ResponseWriter, r *http.Request) {
	var os libocit.OS
	result, _ := ioutil.ReadAll(r.Body)
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

func DeleteOS(w http.ResponseWriter, r *http.Request) {
	Distribution := r.URL.Query().Get("Distribution")
	lock.Lock()
	delete(store, Distribution)
	lock.Unlock()
	w.WriteHeader(http.StatusOK)
}

// Will use DB in the future, (mongodb for example)
func init_db(filename string) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Println(filename, err)
		return
	}

	type _Servers struct {
		Servers []libocit.OS
	}
	var _servers _Servers

	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	json.Unmarshal([]byte(buf.String()), &_servers)

	for index := 0; index < len(_servers.Servers); index++ {
		os := _servers.Servers[index]
		os.Status = "free"
		store[os.ID] = &os
	}

	log.Println(store)
}

var store = map[string]*libocit.OS{}

var lock = sync.RWMutex{}

func main() {
	var config TestServerConfig

	config_content := libocit.ReadFile("./testserver.conf")
	json.Unmarshal([]byte(config_content), &config)

	var port string
	port = fmt.Sprintf(":%d", config.Port)
	//	pub_debug = config.Debug
	init_db(config.ServerListFile)

	mux := routes.New()
	mux.Get("/os", GetOS)
	mux.Post("/os", PostOS)
	mux.Get("/report/:ID", GetReport)
	mux.Post("/casefile/:ID", ReceiveFile)
	mux.Post("/deploy", DeployCommand)
	http.Handle("/", mux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
