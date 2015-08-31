package main

import (
	"encoding/json"
	"github.com/huawei-openlab/oct-engine/lib/libocit"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
)

type TCFile struct {
	Name    string
	Tag     string
	Content string
}

type HostDetail struct {
	Object       string
	Distribution string
	Resource     string
	Containers   string
	Command      string
}
type ContainerDetail struct {
	Class        string
	Distribution string
	ConfigFile   string
	ConfigTag    string
}

type CaseBody struct {
	HostDetails      []HostDetail
	ContainerDetails []ContainerDetail
	TestCase         TCFile
	Files            []TCFile
	Collects         []TCFile
}

func generate_resource(ts_demo libocit.TestCase) (cb CaseBody) {
	for index := 0; index < len(ts_demo.Deploys); index++ {
		var hd HostDetail
		deploy := ts_demo.Deploys[index]
		for r_index := 0; r_index < len(ts_demo.Requires); r_index++ {
			req := ts_demo.Requires[r_index]
			if req.Class == deploy.Class {
				hd.Distribution = req.Distribution + req.Version
				hd.Resource = "CPU " + strconv.Itoa(req.Resource.CPU) + ", Memory " + req.Resource.Memory + ", Disk " + req.Resource.Disk
				break
			}
		}
		hd.Containers = ""
		for c_index := 0; c_index < len(deploy.Containers); c_index++ {
			if len(hd.Containers) > 1 {
				hd.Containers += ", "
			}
			hd.Containers += deploy.Containers[c_index].Object + "(" + deploy.Containers[c_index].Class + ")"

		}
		hd.Object = deploy.Object
		hd.Command = deploy.Cmd
		cb.HostDetails = append(cb.HostDetails, hd)
	}

	for index := 0; index < len(ts_demo.Requires); index++ {
		req := ts_demo.Requires[index]
		if req.Type != "container" {
			continue
		}
		var cd ContainerDetail
		cd.Class = req.Class
		cd.Distribution = req.Distribution + req.Version
		if len(req.Files) > 0 {
			for f_index := 0; f_index < len(req.Files); f_index++ {
				//FIXME: Just want to use container generating file : 'Dockerfile' now
				if path.Base(req.Files[f_index]) == "Dockerfile" {
					cd.ConfigFile = path.Base(req.Files[f_index])
					cd.ConfigTag = strings.ToLower(path.Base(req.Files[f_index]))
					break
				}
			}
		}
		cb.ContainerDetails = append(cb.ContainerDetails, cd)
	}
	return cb
}
func main() {
	var tc libocit.TestCase

	//Head info
	content := libocit.ReadFile("template.json")
	json.Unmarshal([]byte(content), &tc)

	tmpl, err := template.ParseFiles("./template/head.md")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, tc)
	if err != nil {
		panic(err)
	}

	//Summary and Files info
	cb := generate_resource(tc)
	cb.TestCase.Name = "TestCase"
	cb.TestCase.Tag = "testCase"
	cb.TestCase.Content = content

	for index := 0; index < len(tc.Requires); index++ {
		req := tc.Requires[index]
		if req.Type != "container" {
			continue
		}
		if len(req.Files) > 0 {
			for f_index := 0; f_index < len(req.Files); f_index++ {
				//FIXME: Just want to use container generating file : 'Dockerfile' now
				if path.Base(req.Files[f_index]) == "Dockerfile" {
					var c_file TCFile
					c_file.Name = path.Base(req.Files[f_index])
					c_file.Tag = strings.ToLower(path.Base(req.Files[f_index]))
					c_file.Content = libocit.ReadFile(req.Files[f_index])
					cb.Files = append(cb.Files, c_file)
					break
				}
			}
		}
	}

	for index := 0; index < len(tc.Collects); index++ {
		for f_index := 0; f_index < len(tc.Collects[index].Files); f_index++ {
			var c_file TCFile
			file_name := tc.Collects[index].Files[f_index]
			c_file.Name = path.Base(file_name)
			//FIXME: How to convert file to a 'http#' recognized link?
			c_file.Tag = strings.ToLower(c_file.Name)
			c_file.Tag = strings.Replace(c_file.Tag, ".json", "", 1)
			c_file.Content = libocit.ReadFile(file_name)
			cb.Collects = append(cb.Collects, c_file)
		}

	}
	tmpl, err = template.ParseFiles("template/body.md")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, cb)
	if err != nil {
		panic(err)
	}
}
