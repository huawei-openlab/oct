package casevalidator

import (
	"../libocit"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type ValidatorMessage struct {
	//error; warning
	Type string
	Data string
}

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

func FindCaseJson(baseDir string, caseName string) (json_file string, json_dir string) {
	_, err := os.Stat(path.Join(baseDir, "config.json"))
	if err == nil {
		return path.Join(baseDir, "config.json"), baseDir
	}

	files_info, _ := ioutil.ReadDir(baseDir)
	for _, file := range files_info {
		if file.IsDir() {
			sub_json_file, sub_json_dir := FindCaseJson(path.Join(baseDir, file.Name()), caseName)
			if len(sub_json_dir) > 1 {
				return sub_json_file, sub_json_dir
			}
		} else {
			if len(caseName) > 0 {
				if caseName == file.Name() {
					json_file = path.Join(baseDir, file.Name())
					return json_file, baseDir
				}
			} else {
				fileSuffix := path.Ext(file.Name())
				if fileSuffix == ".json" {
					_, err := os.Stat(path.Join(baseDir, "source"))
					if err != nil {
						continue
					} else {
						//  ./casename.json, ./source/
						json_file = path.Join(baseDir, file.Name())
						return json_file, baseDir
					}
				}
			}
		}
	}
	return json_file, json_dir
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

func checkClass(tc libocit.TestCase) (messages []ValidatorMessage) {
	var class_store map[string]bool
	class_store = make(map[string]bool)
	var object_store map[string]string
	object_store = make(map[string]string)

	for index := 0; index < len(tc.Requires); index++ {
		req := tc.Requires[index]
		if len(req.Class) < 1 {
			var msg ValidatorMessage
			msg.Type = "error"
			msg.Data = "No 'Class' in one of 'Requires' session"
			messages = append(messages, msg)
		} else {
			class_store[req.Class] = false
		}
	}

	for index := 0; index < len(tc.Deploys); index++ {
		deploy := tc.Deploys[index]
		if len(deploy.Class) < 1 {
			var msg ValidatorMessage
			msg.Type = "error"
			msg.Data = "No 'Class' in 'Deploys/" + deploy.Object + "' session"
			messages = append(messages, msg)
		} else {
			if _, ok := class_store[deploy.Class]; ok {
				class_store[deploy.Class] = true
				object_store[deploy.Object] = deploy.Class
			} else {
				var msg ValidatorMessage
				msg.Type = "error"
				msg.Data = "The 'Class' " + deploy.Class + " in 'Deploys/" + deploy.Object + "' is not defined"
				messages = append(messages, msg)
			}
		}
	}

	for index := 0; index < len(tc.Run); index++ {
		run := tc.Run[index]
		if len(run.Object) < 1 {
			var msg ValidatorMessage
			msg.Type = "error"
			msg.Data = "No 'Object' in one of the 'Run' session"
			messages = append(messages, msg)
		} else {
			if _, ok := object_store[run.Object]; !ok {
				var msg ValidatorMessage
				msg.Type = "error"
				msg.Data = "The 'Object' " + run.Object + "in one 'Run' session is never deployed"
				messages = append(messages, msg)
			}
		}
	}

	for index := 0; index < len(tc.Collects); index++ {
		collect := tc.Collects[index]
		if len(collect.Object) < 1 {
			var msg ValidatorMessage
			msg.Type = "error"
			msg.Data = "No 'Object' in one of the 'Collect' session"
			messages = append(messages, msg)
		} else {
			if _, ok := object_store[collect.Object]; !ok {
				var msg ValidatorMessage
				msg.Type = "error"
				msg.Data = "The 'Object' " + collect.Object + "in one 'Collect' session is never deployed"
				messages = append(messages, msg)
			}
		}
	}

	for key, val := range class_store {
		if !val {
			var msg ValidatorMessage
			msg.Type = "warning"
			msg.Data = "The 'Class/" + key + "' is never used"
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

func ValidateByFile(caseFile string) {
}

func ValidateByDir(caseDir string, caseName string) (warning_msg []ValidatorMessage, err_msg []ValidatorMessage) {
	var tc libocit.TestCase
	json_file, json_dir := FindCaseJson(caseDir, caseName)

	content := libocit.ReadFile(json_file)

	json.Unmarshal([]byte(content), &tc)

	prop_msgs := checkProp(tc)
	for index := 0; index < len(prop_msgs); index++ {
		if prop_msgs[index].Type == "warning" {
			warning_msg = append(warning_msg, prop_msgs[index])
		} else if prop_msgs[index].Type == "error" {
			err_msg = append(err_msg, prop_msgs[index])
		}
	}

	file_msgs := checkFile(tc, json_dir)
	for index := 0; index < len(file_msgs); index++ {
		if file_msgs[index].Type == "warning" {
			warning_msg = append(warning_msg, file_msgs[index])
		} else if file_msgs[index].Type == "error" {
			err_msg = append(err_msg, file_msgs[index])
		}
	}

	class_msgs := checkClass(tc)
	for index := 0; index < len(class_msgs); index++ {
		if class_msgs[index].Type == "warning" {
			warning_msg = append(warning_msg, class_msgs[index])
		} else if class_msgs[index].Type == "error" {
			err_msg = append(err_msg, class_msgs[index])
		}
	}

	return warning_msg, err_msg
}

/*
func main() {
	var caseDir = flag.String("d", "", "input the case dir")
	var caseFile = flag.String("f", "", "input the file url, case.tar.gz")
	var caseName = flag.String("n", "", "input the 'case name' in the case dir, if there were multiply cases in the case dir. You can use this with -d and -f.")
	var caseID = flag.String("id", "", "input the 'case id' provided by 'Test Case server', please make sure the the tcserver is running.")
	flag.Parse()

	var warning_msg []ValidatorMessage
	var err_msg []ValidatorMessage
	if len(*caseID) > 0 {
		validateByCaseID(*caseID)
	} else if len(*caseFile) > 0 {
		validateByFile(*caseFile)
	} else if len(*caseDir) > 0 {
		warning_msg, err_msg = validateByDir(*caseDir, *caseName)
	} else {
		fmt.Println("Please input the test case")
		return
	}
	if len(err_msg) > 0 {
		fmt.Printf("The case is invalid, there are %d error(errors) and %d warning(warnings)", len(err_msg), len(warning_msg))
		fmt.Println("Please see the details:")
		fmt.Println(err_msg)
		fmt.Println(warning_msg)
	} else if len(warning_msg) > 0 {
		fmt.Printf("The case is OK, but there are %d warning(warnings)", len(warning_msg))
		fmt.Println("Please see the details:")
		fmt.Println(warning_msg)
	} else {
		fmt.Println("Good case.")
	}
}

*/
