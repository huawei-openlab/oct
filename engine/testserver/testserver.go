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
	"path"
	"strconv"
	"strings"
	"sync"
)

type TestServerConfig struct {
	Port           int
	ServerListFile string
	CacheDir       string
	Debug          bool
}

//TODO: not done yet
func GetResult(w http.ResponseWriter, r *http.Request) {
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

//List all the hostOS status
func GetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":ID")
	task := task_store[id]

	body, _ := json.Marshal(*task)
	w.Write([]byte(body))
}

//Set the hostOS status
func SetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":ID")

	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var testStatus libocit.TestingStatus
	json.Unmarshal([]byte(result), &testStatus)

	task := *(task_store[id])

	var ret libocit.HttpRet
	ret.Status = "Failed"
	ret.Message = "Cannot find the related Object " + testStatus.Object
	for index := 0; index < len(task.OSList); index++ {
		if task.OSList[index].Object == testStatus.Object {
			task.OSList[index].Status = testStatus.Status
			ret.Status = "OK"
			ret.Message = "The status is changed."
			break

		}
	}
	ret_string, _ := json.Marshal(ret)
	w.Write([]byte(ret_string))
}

//This is the main task running function
func RunTask(taskID string) {

	var testCommand libocit.TestingCommand
	task := *(task_store[taskID])

	testCommand.ID = taskID
	testCommand.Command = "deploy"

	for index := 0; index < len(task.OSList); index++ {
		//FIXME: should untar the task, compare the resource and send to diffrent OS
		// Here cause it is too late for me to back home..
		testCommand.Object = "hostA"

		fakeIDForHack := "0002"
		os := *(store[fakeIDForHack])
		post_url := "http://" + os.IP + ":9001/command"

		content, _ := json.Marshal(testCommand)
		libocit.SendCommand(post_url, []byte(content))
	}

	testCommand.Command = "run"

	for index := 0; index < len(task.OSList); index++ {
		//FIXME: should untar the task, compare the resource and send to diffrent OS
		// Here cause it is too late for me to back home..
		testCommand.Object = "hostA"

		fakeIDForHack := "0002"
		os := *(store[fakeIDForHack])
		post_url := "http://" + os.IP + ":9001/command"

		content, _ := json.Marshal(testCommand)
		libocit.SendCommand(post_url, []byte(content))
	}
}

func ReceiveTask(w http.ResponseWriter, r *http.Request) {
	//handle_name:  taskID.tar.gz
	real_url, params := libocit.ReceiveFile(w, r, pub_config.CacheDir)

	handle_name := path.Base(real_url)
	taskID := strings.Replace(handle_name, ".tar.gz", "", 1)

	fmt.Println("params id", params["id"], "  ", handle_name)

	var task libocit.Task
	task.ID = taskID
	fakeIDForHack := "0002"

	//FIXME: should untar the task, compare the resource and send to diffrent OS
	// Here cause it is too late for me to back home..
	os := *(store[fakeIDForHack])
	task.OSList = append(task.OSList, os)
	post_url := "http://" + os.IP + ":9001/task"
	fmt.Println("Receive and send file ", real_url)
	w.Write([]byte("receive and send file"))
	task_store[taskID] = &task

	//FIXME: it is better to send the related the file to the certain host OS
	libocit.SendFile(post_url, real_url, params)

	//FIXME: if there were not enough resource ,return error
	var ret libocit.HttpRet
	ret.Status = "OK"
	ret.Message = "success in receiving task files"

	ret_string, _ := json.Marshal(ret)
	w.Write([]byte(ret_string))

	RunTask(taskID)
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
	var ret libocit.HttpRet

	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if pub_config.Debug {
		fmt.Println(string(result))
	}
	json.Unmarshal([]byte(result), &os)
	if os.Distribution == "" {
		ret.Status = "Failed"
		ret.Message = "os distribution required"
	} else if os.Version == "" {
		ret.Status = "Failed"
		ret.Message = "os version required"
	} else if os.Arch == "" {
		ret.Status = "Failed"
		ret.Message = "os arch required"
	} else {
		lock.Lock()
		id := libocit.MD5(string(result))
		if _, ok := store[id]; ok {
			ret.Status = "Failed"
			ret.Message = "this os is already exist"
		} else {
			store[id] = &os
			ret.Status = "OK"
			ret.Message = "Success in adding the os"
		}
		lock.Unlock()
	}
	ret_body, _ := json.Marshal(ret)
	w.Write([]byte(ret_body))
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
		content, _ := json.Marshal(os)
		id := libocit.MD5(string(content))
		store[id] = &os
	}

	log.Println(store)
}

var store = map[string]*libocit.OS{}
var task_store = map[string]*libocit.Task{}

var lock = sync.RWMutex{}
var pub_config TestServerConfig

func main() {

	config_content := libocit.ReadFile("./testserver.conf")
	json.Unmarshal([]byte(config_content), &pub_config)

	var port string
	port = fmt.Sprintf(":%d", pub_config.Port)
	init_db(pub_config.ServerListFile)

	mux := routes.New()
	//TODO: following one is are not done yet
	mux.Get("/os", GetOS)

	mux.Post("/os", PostOS)
	mux.Get("/:ID/status", GetStatus)
	mux.Post("/:ID/status", SetStatus)
	mux.Post("/task", ReceiveTask)

	mux.Get("/:ID/result", GetResult)

	http.Handle("/", mux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
