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
	"strconv"
	"strings"
)

type TCServerConf struct {
	GitRepo  string
	CaseDir  string
	Group    []string
	CacheDir string
	Port     int
}

type MetaUnit struct {
	ID     string
	Group  string
	Name   string
	Status string
	//0 means not tested
	TestedTime       int64
	LastModifiedTime int64
}

var store = map[string]*MetaUnit{}
var pub_config TCServerConf

func RefreshRepo() {
	var cmd string
	libocit.PreparePath(pub_config.CacheDir, "")
	repo_name := strings.Replace(path.Base(pub_config.GitRepo), ".git", "", 1)
	//FIXME: better way? using github go lib in the future
	git_check_url := path.Join(pub_config.CacheDir, repo_name, ".git/config")
	_, err := os.Stat(git_check_url)

	if err != nil {
		cmd = "cd " + pub_config.CacheDir + " ; git clone " + pub_config.GitRepo
	} else {
		cmd = "cd " + path.Join(pub_config.CacheDir, repo_name) + " ; git pull"
	}

	fmt.Println("Refresh by using ", cmd)
	c := exec.Command("/bin/sh", "-c", cmd)
	c.Run()
	fmt.Println("Refresh done")
}

//TODO, since case validation is not implemented now, return 0 means the case is invalid
func LastModified(case_dir string) (last_modified int64) {
	config_url := path.Join(case_dir, "config.json")
	fi, err := os.Stat(config_url)
	last_modified = 0
	if err != nil {
		return last_modified
	} else {
		last_modified = fi.ModTime().Unix()
	}
	files, _ := ioutil.ReadDir(path.Join(config_url, "source"))
	for _, file := range files {
		if file.IsDir() {
			//This case format is not suggested
			continue
		} else {
			if last_modified < file.ModTime().Unix() {
				last_modified = file.ModTime().Unix()
			}
		}
	}
	return last_modified
}

func LoadDB() {
	RefreshRepo()

	for g_index := 0; g_index < len(pub_config.Group); g_index++ {
		repo_name := strings.Replace(path.Base(pub_config.GitRepo), ".git", "", 1)
		group_dir := path.Join(pub_config.CacheDir, repo_name, pub_config.CaseDir, pub_config.Group[g_index])
		files, _ := ioutil.ReadDir(group_dir)
		for _, file := range files {
			if file.IsDir() {
				//TODO, Qilin is working on case validation work, here we should check it!
				//	or we can check it in case push phase
				last_modified := LastModified(path.Join(group_dir, file.Name()))
				if last_modified == 0 {
					continue
				}

				store_md := libocit.MD5(path.Join(pub_config.Group[g_index], file.Name()))
				if v, ok := store[store_md]; ok {
					//Happen when we refresh the repo
					(*v).LastModifiedTime = last_modified
					fi, err := os.Stat(path.Join(group_dir, file.Name(), "report.md"))
					if err != nil {
						(*v).TestedTime = 0
					} else {
						(*v).TestedTime = fi.ModTime().Unix()
					}
					if (*v).LastModifiedTime > (*v).TestedTime {
						(*v).Status = "idle"
					} else {
						(*v).Status = "tested"
					}
				} else {
					var meta MetaUnit
					meta.ID = store_md
					meta.Group = pub_config.Group[g_index]
					meta.Name = file.Name()
					fi, err := os.Stat(path.Join(group_dir, file.Name(), "report.md"))
					if err != nil {
						meta.TestedTime = 0
					} else {
						meta.TestedTime = fi.ModTime().Unix()
					}
					meta.LastModifiedTime = last_modified
					if meta.LastModifiedTime > meta.TestedTime {
						meta.Status = "idle"
					} else {
						meta.Status = "tested"
					}
					store[store_md] = &meta
				}
			}
		}
	}
}

func ListCases(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("Status")
	page_string := r.URL.Query().Get("Page")
	page, err := strconv.Atoi(page_string)
	if err != nil {
		page = 0
	}
	page_size_string := r.URL.Query().Get("PageSize")
	page_size, err := strconv.Atoi(page_size_string)
	if err != nil {
		page_size = 10
	}

	var case_list []MetaUnit
	cur_num := 0
	for _, tc := range store {
		if status != "" {
			if status != tc.Status {
				continue
			}
		}
		cur_num += 1
		if (cur_num >= page*page_size) && (cur_num < (page+1)*page_size) {
			case_list = append(case_list, *tc)
		}

	}

	case_string, err := json.Marshal(case_list)
	if err != nil {
		w.Write([]byte("[]"))
	} else {
		w.Write([]byte(case_string))
	}

}

func GetCase(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":ID")
	meta := store[id]
	repo_name := strings.Replace(path.Base(pub_config.GitRepo), ".git", "", 1)
	case_dir := path.Join(pub_config.CacheDir, repo_name, pub_config.CaseDir, meta.Group, meta.Name)
	tar_url := libocit.TarDir(case_dir)

	file, err := os.Open(tar_url)
	defer file.Close()
	if err != nil {
		//FIXME: add to head
		w.Write([]byte("Cannot open the file: " + tar_url))
		return
	}

	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	//TODO: write head, filename and the etc
	w.Write([]byte(buf.String()))
}

func GetCaseReport(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":ID")
	meta := store[id]
	repo_name := strings.Replace(path.Base(pub_config.GitRepo), ".git", "", 1)
	report_url := path.Join(pub_config.CacheDir, repo_name, pub_config.CaseDir, meta.Group, meta.Name, "report.md")

	_, err := os.Stat(report_url)
	if err != nil {
		//FIXME: 404 error head
		w.Write([]byte("Cannot find the report"))
		return
	}
	content := libocit.ReadFile(report_url)
	w.Write([]byte(content))
}

func RefreshCases(w http.ResponseWriter, r *http.Request) {
	RefreshRepo()
	var ret libocit.HttpRet
	ret.Status = "OK"
	ret_string, _ := json.Marshal(ret)
	w.Write([]byte(ret_string))
}

func main() {
	content := libocit.ReadFile("./tcserver.conf")
	json.Unmarshal([]byte(content), &pub_config)
	LoadDB()

	port := fmt.Sprintf(":%d", pub_config.Port)
	fmt.Println("Listen to port ", port)
	mux := routes.New()
	mux.Get("/case", ListCases)
	mux.Post("/case", RefreshCases)
	mux.Get("/case/:ID", GetCase)
	mux.Get("/case/:ID/report", GetCaseReport)
	http.Handle("/", mux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
