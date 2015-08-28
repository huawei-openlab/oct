package main

import (
	"../../lib/libocit"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

var pub_casedir string

func generate_head(ts_demo libocit.TestCase) (content string) {
	content = "## " + ts_demo.Name + "-" + ts_demo.Version + "\n"
	content += "[Test Case](#testcase) " + ts_demo.Description + "\n\n"
	content += "```\n"
	content += "Owner: " + ts_demo.Owner + "\n"
	content += "License: " + ts_demo.License + "\n"
	content += "Group: " + ts_demo.Group + "\n"
	content += "```\n\n"

	return content
}

func generate_resource(ts_demo libocit.TestCase) (content string) {
	num := len(ts_demo.Deploys)
	content = "The case has " + strconv.Itoa(num) + " host operation system(s):\n\n"

	for index := 0; index < num; index++ {
		deploy := ts_demo.Deploys[index]
		content += "'" + deploy.Object + "' has " + strconv.Itoa(len(deploy.Containers)) + " container(s) deployed.\n\n"
	}
	content += "The defailed information is listed as below:\n\n"

	content += "| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |\n"
	content += "| -------| ------ | --------- | -------- | --------|\n"
	for index := 0; index < num; index++ {
		deploy := ts_demo.Deploys[index]
		var version string
		var resource string
		var container string
		for r_index := 0; r_index < len(ts_demo.Requires); r_index++ {
			req := ts_demo.Requires[r_index]
			if req.Class == deploy.Class {
				version = req.Distribution + req.Version
				resource = "CPU " + strconv.Itoa(req.Resource.CPU) + ", Memory " + req.Resource.Memory + ", Disk " + req.Resource.Disk
				break
			}
		}
		container = ""
		for c_index := 0; c_index < len(deploy.Containers); c_index++ {
			if len(container) > 1 {
				container += ", "
			}
			container += deploy.Containers[c_index].Object + "(" + deploy.Containers[c_index].Class + ")"

		}
		content += "|" + deploy.Object + "|" + version + "|" + resource + "|" + container + "|\"" + deploy.Cmd + "\"|\n"
	}

	content += "\nThe defailed information of each container type is listed as below:\n\n"
	content += "| *Container Type* | *Distribution* | *Container File* |\n"
	content += "| -------| ------ | ------- |\n"
	for index := 0; index < len(ts_demo.Requires); index++ {
		req := ts_demo.Requires[index]
		if req.Type != "container" {
			continue
		}
		content += "|" + req.Class + "|" + req.Distribution + req.Version + "|" + "[Dockerfile](#dockerfile) |\n"
	}
	return content
}

func generate_result_link(ts_demo libocit.TestCase) (content string) {
	content = "\nAfter running the `Command` in each OS and container, we get two results.\n\n"

	for index := 0; index < len(ts_demo.Collects); index++ {
		collect := ts_demo.Collects[index]
		for f_index := 0; f_index < len(collect.Files); f_index++ {
			basename := path.Base(collect.Files[f_index])
			basename_without_json := strings.Replace(basename, ".json", "", 1)
			content += "* [" + basename + "](#" + basename_without_json + ") \n"
		}
	}

	return content
}

func generate_result_file(case_dir string, ts_demo libocit.TestCase) (content string) {
	for index := 0; index < len(ts_demo.Collects); index++ {
		collect := ts_demo.Collects[index]
		for f_index := 0; f_index < len(collect.Files); f_index++ {

			basename := path.Base(collect.Files[f_index])
			basename_without_json := strings.Replace(basename, ".json", "", 1)
			content += "\n###" + basename_without_json + "\n"
			content += "```\n"

			_, err := os.Stat(collect.Files[f_index])
			if err != nil {
				content += libocit.ReadFile(path.Join(case_dir, collect.Files[f_index]))
			} else {
				content += libocit.ReadFile(collect.Files[f_index])
			}
			content += "\n```\n\n"
		}
	}
	return content
}

func main() {
	var ts_demo libocit.TestCase
	var case_file string

	arg_num := len(os.Args)
	if arg_num < 2 {
		fmt.Println("Please input the testcase dir")
		return
	} else {
		case_file = os.Args[1]
	}

	pub_casedir = path.Dir(case_file)
	test_json_str := libocit.ReadFile(case_file)
	json.Unmarshal([]byte(test_json_str), &ts_demo)

	content := generate_head(ts_demo)
	content += generate_resource(ts_demo)

	content += generate_result_link(ts_demo)

	content += "\n\n###TestCase\n"
	content += "```\n"
	content += test_json_str
	content += "```\n\n"

	content += "\n\n###Dockerfile\n"
	content += "```\n"
	content += libocit.ReadFile(path.Join(pub_casedir, "./source/Dockerfile"))
	content += "```\n\n"

	content += generate_result_file(pub_casedir, ts_demo)

	fmt.Println(content)
}
