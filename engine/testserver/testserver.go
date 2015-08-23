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
//TODO, In the async testing, when all the OS get the same status, continue to the next
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

//TODO make it async 
//This is the main task running function
func RunTask(taskID string) {

	var testCommand libocit.TestingCommand
	task := *(task_store[taskID])

	//Deploy
	testCommand.ID = taskID
	testCommand.Status = "deploy"
	for index := 0; index < len(task.TC.Deploys); index++ {
		id := task.TC.Deploys[index].ID
		hostOS := *(store[id])
		//TODO: 9001 will be replaced with OS.Port
		post_url := "http://" + hostOS.IP + ":9001/command"

		testCommand.Command = task.TC.Deploys[index].Cmd
		content, _ := json.Marshal(testCommand)
		if pub_config.Debug {
			fmt.Println("Command ", string(content), "  to ", post_url)
		}
		libocit.SendCommand(post_url, []byte(content))

		//TODO: add the container command 
	}

	//Run

	testCommand.Status = "run"
	for index := 0; index < len(task.TC.Run); index++ {
		id := task.TC.Run[index].ID
		hostOS := *(store[id])
		//TODO: 9001 will be replaced with OS.Port
		post_url := "http://" + hostOS.IP + ":9001/command"
		testCommand.Command = task.TC.Run[index].Cmd
		content, _ := json.Marshal(testCommand)
		if pub_config.Debug {
			fmt.Println("Command ", string(content), "  to ", post_url)
		}
		libocit.SendCommand(post_url, []byte(content))

		//TODO: add the container command 
	}

	//Collect
	for index := 0; index < len(task.TC.Collects); index++ {
		collect := task.TC.Collects[index]
		id := collect.ID
		hostOS := *(store[id])
		for f_index := 0; f_index < len(collect.Files); f_index++ {
			//TODO: 9001 will be replaced with OS.Port
			get_url := "http://" + hostOS.IP + ":9001/result" + "?File=" + collect.Files[f_index] + "&ID=" + taskID
			if pub_config.Debug {
				fmt.Println("Get collect url : ", get_url)
			}

			resp, err := http.Get(get_url)
			if err != nil {
				fmt.Println("Error ", err)
				continue
			}
			defer resp.Body.Close()
			resp_body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}
			fmt.Println(resp.Status)
			fmt.Println(string(resp_body))
			real_url := libocit.PreparePath(path.Join(pub_config.CacheDir, taskID), collect.Files[f_index])
			f, err := os.Create(real_url)
			defer f.Close()
			f.Write(resp_body)
			f.Sync()
		}

	}
}

func AllocateOS(task libocit.Task) (success bool) {
	for index := 0; index < len(task.TC.Requires); index++ {
		req := task.TC.Requires[index]
		if req.Type == "os" {
			var os_query libocit.OS
			//TODO, now only check these options
			os_query.Distribution = req.Distribution
			os_query.Version = req.Version
			id := GetAvaliableResource(os_query)
			if len(id) < 1 {
				success = false
				return success
			}
			task.OSList = append(task.OSList, *(store[id]))
			for d_index := 0; d_index < len(task.TC.Deploys); d_index++ {
				deploy := task.TC.Deploys[d_index]
				if deploy.Class == req.Class {
					task.TC.Deploys[d_index].ID = id
					for r_index := 0; r_index < len(task.TC.Run); r_index++ {
						run := task.TC.Run[r_index]
						if deploy.Object == run.Object {
							task.TC.Run[r_index].ID = id
						}

					}
					for c_index := 0; c_index < len(task.TC.Collects); c_index++ {
						collect := task.TC.Collects[c_index]
						if deploy.Object == collect.Object {
							task.TC.Collects[c_index].ID = id
						}

					}

				}

			}
			if pub_config.Debug {
				ret_string, _ := json.Marshal(task.TC)
				fmt.Println("get --- id ---- ", string(ret_string))
			}
		}
	}
	success = true
	return success
}

func ReceiveTask(w http.ResponseWriter, r *http.Request) {
	real_url, params := libocit.ReceiveFile(w, r, pub_config.CacheDir)

	taskID := params["id"]

	var task libocit.Task
	task.ID = taskID

	// for example, we have taskID.tar.gz
	//  untar it, the test case will be put into taskID/config.json
	// Should always use 'config.json'
	// content := libocit.ReadTar(real_url, "config.json", "")
	content := libocit.ReadCaseFromTar(real_url)
	if pub_config.Debug {
		fmt.Println(content)
	}
	json.Unmarshal([]byte(content), &task.TC)
	success := AllocateOS(task)
	if success == false {
		var ret libocit.HttpRet
		ret.Status = "Failed"
		ret.Message = "Donnot have the required operation systems"

		ret_string, _ := json.Marshal(ret)
		if pub_config.Debug {
			fmt.Println(string(ret_string))
		}
		w.Write([]byte(ret_string))
		return
	}

	task_store[taskID] = &task
	var success_count int
	success_count = 0
	for index := 0; index < len(task.TC.Deploys); index++ {
		deploy := task.TC.Deploys[index]
		os := *(store[deploy.ID])
		post_url := "http://" + os.IP + ":9001/task"
		fmt.Println("Receive and send file ", real_url, " to  ", post_url)

		//FIXME: it is better to send the related the file to the certain host OS
		ret := libocit.SendFile(post_url, real_url, params)
		if ret.Status == "OK" {
			success_count += 1
		} else {
			if pub_config.Debug {
				fmt.Println(ret)
			}
		}

	}

	var ret libocit.HttpRet
	if success_count == len(task.TC.Deploys) {
		ret.Status = "OK"
		ret.Message = "success in receiving task files"
	} else {
		ret.Status = "Failed"
		ret.Message = "Some testing files were not send successfully"
	}

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
//TODO: different ID even same Class, add the ticket?
func GetAvaliableResource(os_query libocit.OS) (ID string) {
	for _, os := range store {
		if len(os_query.Distribution) > 1 {
			if strings.EqualFold(os_query.Distribution, (*os).Distribution) == false {
				continue
			}
		}
		if len(os_query.Version) > 1 {
			//TODO do not check the version for now
			//			if os_query.Version != (*os).Version {
			//				continue
			//			}
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

func GetResourceList(os_query libocit.OS) (ids []string) {
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
		ids = append(ids, (*os).ID)
	}
	return ids
}

func GetOS(w http.ResponseWriter, r *http.Request) {
	var os_query libocit.OS
	os_query = GetOSQuery(r)

	ids := GetResourceList(os_query)

	var ret libocit.HttpRet
	if len(ids) < 1 {
		ret.Status = "Failed"
		ret.Message = "Cannot find the avaliable OS"
	} else {
		ret.Status = "OK"
		ret.Message = "Find the avaliable OS"
		var oss []libocit.OS
		for index := 0; index < len(ids); index++ {
			oss = append(oss, *(store[ids[index]]))
		}

		data, _ := json.Marshal(oss)
		ret.Data = string(data)
	}

	body, _ := json.Marshal(ret)
	w.Write([]byte(body))
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
			os.ID = id
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
		os.ID = id
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
