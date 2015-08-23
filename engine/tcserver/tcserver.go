package main

import (
	"../lib/libocit"
	"../lib/casevalidator"
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
	CaseFolderName  string
	Groups    []Group
	CacheDir string
	Port     int
}

type Group struct {
	Name string
	LibFolderName string
}

type MetaUnit struct {
	ID     string
	Name   string
	GroupDir string
	LibFolderName string
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

func LastModified(case_dir string) (last_modified int64) {
	last_modified = 0
	files, _ := ioutil.ReadDir(case_dir)
	for _, file := range files {
		if file.IsDir() {
			sub_lm := LastModified(path.Join(case_dir, file.Name()))
			if last_modified < sub_lm {
				last_modified = sub_lm
			}
		} else {
			if last_modified < file.ModTime().Unix() {
				last_modified = file.ModTime().Unix()
			}
		}
	}
	return last_modified
}

func LoadCase(groupDir string, caseName string, caseLibFolderName string) {
	caseDir := path.Join(groupDir, caseName)
	_, err_msgs := casevalidator.ValidateByDir(caseDir, "")
	if len(err_msgs) == 0 {
				last_modified := LastModified(caseDir)
				store_md := libocit.MD5(caseDir)
				if v, ok := store[store_md]; ok {
					//Happen when we refresh the repo
					(*v).LastModifiedTime = last_modified
					fi, err := os.Stat(path.Join(caseDir, "report.md"))
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
					meta.Name = caseName
					meta.GroupDir = groupDir
					meta.LibFolderName = caseLibFolderName
					fi, err := os.Stat(path.Join(caseDir, "report.md"))
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
	} else {
		fmt.Println("Error in loading case: ", caseDir, " . Skip it")
		return
	}
}

func LoadCaseGroup(groupDir string, libDir string) {
	files, _ := ioutil.ReadDir(groupDir)
	for _, file := range files {
		if file.IsDir() {
			if len(libDir) > 0 {
				if libDir == file.Name() {
					continue
				} else {
					LoadCase(groupDir, file.Name(), libDir)
				}
			} else {
				LoadCase(groupDir, file.Name(), "")
			}
		}
	}
}

func LoadDB() {
	RefreshRepo()

	for g_index := 0; g_index < len(pub_config.Groups); g_index++ {
		repo_name := strings.Replace(path.Base(pub_config.GitRepo), ".git", "", 1)
		group_dir := path.Join(pub_config.CacheDir, repo_name, pub_config.CaseFolderName, pub_config.Groups[g_index].Name)
		LoadCaseGroup(group_dir, pub_config.Groups[g_index].LibFolderName)
/*
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
*/
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
	files := libocit.GetDirFiles(meta.GroupDir, meta.Name)
	if len(meta.LibFolderName) > 0 {
		lib_files := libocit.GetDirFiles(meta.GroupDir, meta.LibFolderName)
		for index := 0; index < len(lib_files); index++ {
			files = append(files, lib_files[index])
		}
	}
	tar_url := libocit.TarFileList(files, meta.GroupDir, meta.Name)

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
	report_url := path.Join(pub_config.CacheDir, repo_name, pub_config.CaseFolderName, meta.GroupDir, meta.Name, "report.md")

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
