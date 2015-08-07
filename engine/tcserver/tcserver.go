package main

import (
	"../lib/libocit"
	"../lib/routes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	//	"os/exec"
	"bytes"
	"path"
)

type TCServerConf struct {
	GitRepo  string
	CaseDir  string
	Group    []string
	CacheDir string
	Metafile string
	Port     int
}

type MetaUnit struct {
	Group string
	Name  string

	//0 means not tested
	TestedTime       int64
	LastModifiedTime int64
}

var store = map[string]*MetaUnit{}
var pub_config TCServerConf

func RefreshRepo(repo string) {
}

func LoadMeta(meta_file string) {
	type _Metas struct {
		Metas []MetaUnit
	}
	var _metas _Metas
	var content string
	libocit.ReadFile(meta_file)
	json.Unmarshal([]byte(content), &_metas)
	metas := _metas.Metas
	for index := 0; index < len(metas); index++ {
		url := path.Join(metas[index].Group, metas[index].Name)
		if v, ok := store[url]; ok {
			fmt.Println("Error in meta file, duplicated testcase record: ", v)
		} else {
			store[url] = &metas[index]
		}
	}
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
	RefreshRepo(pub_config.GitRepo)
	LoadMeta(pub_config.Metafile)

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
					if (*v).LastModifiedTime < last_modified {
						(*v).LastModifiedTime = last_modified
					}
				} else {
					var meta MetaUnit
					meta.Group = pub_config.Group[g_index]
					meta.Name = file.Name()
					meta.TestedTime = 0
					meta.LastModifiedTime = last_modified
					store[store_md] = &meta
				}
			}
		}
	}
}

func ListCases(w http.ResponseWriter, r *http.Request) {
	//TODO: add 'Query' support
	store_string, _ := json.Marshal(store)
	fmt.Println(store_string)
	w.Write([]byte(store_string))

}

func GetCase(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":ID")
	meta := store[id]
	//      meta_string, _ := json.Marshal(meta)
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

func main() {
	content := libocit.ReadFile("./tcserver.conf")
	json.Unmarshal([]byte(content), &pub_config)
	LoadDB()

	port := fmt.Sprintf(":%d", pub_config.Port)
	fmt.Println("Listen to port ", port)
	mux := routes.New()
	mux.Get("/case", ListCases)
	mux.Get("/case/:ID", GetCase)
	//	mux.Post("/refresh", RefreshDB)
	//	mux.Post("/report", AddReport)
	http.Handle("/", mux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
