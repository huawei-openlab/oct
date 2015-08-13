package main

import (
	"../../lib/libocit"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// The case now could be like this:
//  in this case type, we will send all the files to all the hostOS
//   casegroup 
//          |____ casedir
//          |         |___ casename.json
//          |         |___ `source`        (must be `source`)
//          |                  |____ file1
//          |                  |____ ...
//          |                  |____ fileN
//          |                  |____ dir1
//          |                  |____ ...
//          |                  |____ dirN
//          |                 
//          |____  caselibdir
//                    |_____ libfile1
//                    |_____  ....
//                    |_____ libfile2
//
//
// The ideal case should be like this:
//
//   casedir
//        |___  `config.json` (must be `config.json`
//        |___  `source`      (must be `source` dir)
//                  |____ file1
//                  |____  ...
//                  |____ fileN
//                  |____ dir1 with files
//                  |____  ...
//                  |____ dirN with files
//

func FindCaseJson(base_dir string) (json_file string, json_dir string) {
	_, err := os.Stat(path.Join(base_dir, "config.json"))
	if err == nil {
		return path.Join(base_dir, "config.json"), base_dir
	}

	files_info, _ := ioutil.ReadDir(base_dir)
	for _, file := range files_info {
		if file.IsDir() {
			sub_json_file, sub_json_dir := FindCaseJson(path.Join(base_dir, file.Name()))
			if len(sub_json_dir) > 1 {
				return sub_json_file, sub_json_dir
			}
		} else {
			fileSuffix := path.Ext(file.Name())
			if fileSuffix == ".json" {
				_, err := os.Stat(path.Join(base_dir, "source"))
				if err != nil {
					continue
				} else {
					//  ./casename.json, ./source/
					json_file = path.Join(base_dir, file.Name())
					return json_file, base_dir
				}
			}
		}
	}
	return json_file, json_dir
}

type ValidatorMessage struct {
	//error; warning
	Type string
	Data string
}

func checkProp(tc libocit.TestCase) (messages []ValidatorMessage) {
	if len(tc.Name) < 1 {
		var msg ValidatorMessage
		msg.Type = "error"
		msg.Data = "'Name' not found."
		messages = append(messages, msg)
	}
	if len(tc.Version) < 1 {
		var msg ValidatorMessage
		msg.Type = "warning"
		msg.Data = "'Version' not found."
		messages = append(messages, msg)
	}
	if len(tc.License) < 1 {
		var msg ValidatorMessage
		msg.Type = "warning"
		msg.Data = "'License' not found."
		messages = append(messages, msg)
	}
	if len(tc.Group) < 1 {
		var msg ValidatorMessage
		msg.Type = "warning"
		msg.Data = "'Group' not found. Please read the 'Group' defination in OCT"
		messages = append(messages, msg)
	}
	if len(tc.Owner) < 1 {
		var msg ValidatorMessage
		msg.Type = "warning"
		msg.Data = "'Owner' not found."
		messages = append(messages, msg)
	}
	if len(tc.Sources) > 0 {
		var msg ValidatorMessage
		msg.Type = "warning"
		msg.Data = "Don't need to add `Source in this part."
		messages = append(messages, msg)
	}

	if len(tc.Requires) == 0 {
		var msg ValidatorMessage
		msg.Type = "error"
		msg.Data = "No 'Requires' found, we don't know what kind of resource your case need."
		messages = append(messages, msg)
	}

	if len(tc.Deploys) == 0 {
		var msg ValidatorMessage
		msg.Type = "error"
		msg.Data = "No 'Deploys' found, we don't know how to deploy the case."
		messages = append(messages, msg)
	} else {
		for d_index := 0; d_index < len(tc.Deploys); d_index++ {
			deploy := tc.Deploys[d_index]
			if (len(deploy.Cmd) == 0) && (len(deploy.Files) == 0) {
				var msg ValidatorMessage
				msg.Type = "warning"
				msg.Data = "No 'Cmd' and 'Files in 'Deploys/" + deploy.Object + "' found, maybe we can remove the object?"
				messages = append(messages, msg)
			}
		}
	}

	if len(tc.Run) == 0 {
		var msg ValidatorMessage
		msg.Type = "warning"
		msg.Data = "No 'Run' found, if you put the running command to the 'Deploy' session, that is OK."
		messages = append(messages, msg)
	} else {
		for r_index := 0; r_index < len(tc.Run); r_index++ {
			run := tc.Run[r_index]
			if len(run.Files) > 0 {
				var msg ValidatorMessage
				msg.Type = "warning"
				msg.Data = "It is OK to put files in 'Run' session in the " + run.Object + ". But we suggest to move it to 'Deploys' session."
				messages = append(messages, msg)
			}
			if len(run.Cmd) == 0 {
				var msg ValidatorMessage
				msg.Type = "warning"
				msg.Data = "No 'Cmd' in 'Run/" + run.Object + "' session. maybe we can remove the object?"
				messages = append(messages, msg)
			}
		}
	}

	if len(tc.Collects) == 0 {
		var msg ValidatorMessage
		msg.Type = "error"
		msg.Data = "No 'Collects' found, we need the testing result to generate the report."
		messages = append(messages, msg)
	} else {
		haveCollectedFile := false
		for c_index := 0; c_index < len(tc.Collects); c_index++ {
			if len(tc.Collects[c_index].Files) == 0 {
				var msg ValidatorMessage
				msg.Type = "warning"
				msg.Data = "No 'Files' in 'Collects/" + tc.Collects[c_index].Object + "' session, maybe we can remove the object?"
				messages = append(messages, msg)
			} else {
				haveCollectedFile = true
			}
		}
		if haveCollectedFile == false {
			var msg ValidatorMessage
			msg.Type = "error"
			msg.Data = "No 'Files' in 'Collects' found, we need the testing result to generate the report."
			messages = append(messages, msg)
		}
	}
	return messages
}

func checkFile(tc libocit.TestCase, casedir string) (messages []ValidatorMessage) {

	var file_store map[string]string
	file_store = make(map[string]string)

	files := libocit.GetDirFiles(casedir, "source")
	for index := 0; index < len(files); index++ {
		file_store[path.Join(casedir, files[index])] = files[index]
	}

	for index := 0; index < len(tc.Deploys); index++ {
		deploy := tc.Deploys[index]
		for f_index := 0; f_index < len(deploy.Files); f_index++ {
			file := path.Join(casedir, deploy.Files[f_index])
			if _, ok := file_store[file]; ok {
				file_store[file] = ""
			} else {
				var msg ValidatorMessage
				msg.Type = "error"
				msg.Data = "File " + file + " mentioned in 'Deploys/" + deploy.Object + "' part is not exist. Forget to submit the file?"
				messages = append(messages, msg)
			}
		}
	}

	for index := 0; index < len(tc.Run); index++ {
		run := tc.Run[index]
		for f_index := 0; f_index < len(run.Files); f_index++ {
			file := path.Join(casedir, run.Files[f_index])
			if _, ok := file_store[file]; ok {
				file_store[file] = ""
			} else {
				var msg ValidatorMessage
				msg.Type = "error"
				msg.Data = "File " + file + " mentioned in 'Run/" + run.Object + "' part is not exist. Forget to submit the file?"
				messages = append(messages, msg)
			}
		}
	}

	for _, val := range file_store {
		if len(val) > 0 {
			var msg ValidatorMessage
			msg.Type = "warning"
			msg.Data = "File " + val + " in the test case directory is not mentioned in the case file."
			messages = append(messages, msg)
		}
	}
	return messages
}

func validateByCaseID(caseID string) {
}

func validateByFile(caseFile string) {
}

func validateByDir(caseDir string) {
	var tc libocit.TestCase
	json_file, json_dir := FindCaseJson(caseDir)
	content := libocit.ReadFile(json_file)

	json.Unmarshal([]byte(content), &tc)

	fmt.Println(json_dir)

	prop_msgs := checkProp(tc)
	fmt.Println(prop_msgs)

	file_msgs := checkFile(tc, json_dir)
	fmt.Println(file_msgs)
}

func main() {
	var caseDir = flag.String("d", "", "input the case dir")
	var caseFile = flag.String("f", "", "input the file url, case.tar.gz")
	var caseID = flag.String("id", "", "input the 'case id' provided by 'Test Case server', please make sure the the tcserver is running.")
	flag.Parse()

	if len(*caseID) > 0 {
		validateByCaseID(*caseID)
	} else if len(*caseFile) > 0 {
		validateByFile(*caseFile)
	} else if len(*caseDir) > 0 {
		validateByDir(*caseDir)
	} else {
		fmt.Println("Please input the test case")
	}
}
