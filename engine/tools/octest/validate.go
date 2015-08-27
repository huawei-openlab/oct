// Copyright 2015 The oct Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"github.com/codegangsta/cli"
        "github.com/opencontainers/specs"

)

const (
	// Path to config file inside the layout
	ConfigFile = "config.json"
	// Path to rootfs directory inside the layout
	RootfsDir = "rootfs"
)

var (
	ErrNoRootFS   = errors.New("no rootfs found in layout")
	ErrNoConfig = errors.New("no config json file found in layout")
)

func validate(context *cli.Context) {
    args := context.String("config")
    
    if len(args) == 0 {
        args = context.String("layout")
        if len(args) == 0 {
            cli.ShowCommandHelp(context, "validate")
            return
        } else {
           err := validateLayout(args) 
           if err != nil {
				fmt.Printf("%s: invalid image layout: %v\n", args, err)
			} else {
				fmt.Printf("%s: valid image layout\n", args)
			}           
        }
    } else {
               validateConfigFile(args)
    }


}

func validateLayout(path string) error {
        fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error accessing layout: %v", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("given path %q is not a directory", path)
	}
	var flist []string
	var imOK, rfsOK bool
	var im io.Reader
	walkLayout := func(fpath string, fi os.FileInfo, err error) error {
		rpath, err := filepath.Rel(path, fpath)
		if err != nil {
			return err
		}
		switch rpath {
		case ".":
		case ConfigFile:
			im, err = os.Open(fpath)
			if err != nil {
				return err
			}
			imOK = true
		case RootfsDir:
			if !fi.IsDir() {
				return errors.New("rootfs is not a directory")
			}
			rfsOK = true
		default:
			flist = append(flist, rpath)
		}
		return nil
	}
	if err := filepath.Walk(path, walkLayout); err != nil {
		return err
	}
	return checkLayout(imOK, im, rfsOK, flist)
}

func checkLayout(imOK bool, im io.Reader, rfsOK bool, files []string) error {
	defer func() {
		if rc, ok := im.(io.Closer); ok {
			rc.Close()
		}
	}()
	if !imOK {
		return ErrNoConfig
	}
	if !rfsOK {
		return ErrNoRootFS
	}
	_, err := ioutil.ReadAll(im)
	if err != nil {
		return fmt.Errorf("error reading the layout: %v", err)
	}
		
	for _, f := range files {
		if !strings.HasPrefix(f, "rootfs") {
			return fmt.Errorf("unrecognized file path in layout: %q", f)
		}
	}
	return nil
}


func validateConfigFile(path string) {
	var sp specs.Spec
	content, err := ReadFile(path)
	if err != nil {
		return
	}
	json.Unmarshal([]byte(content), &sp)
	var secret interface{} = sp
	value := reflect.ValueOf(secret)

	var err_msg []string
	ok, err_msg := validateStruct(value, reflect.TypeOf(secret).Name(), err_msg)

	if ok == false {
		fmt.Println("The configuration is incomplete, see the details: \n")
		for index := 0; index < len(err_msg); index++ {
			fmt.Println(err_msg[index])
		}
	} else {
		fmt.Println("The configuration is Good")

	}


}

func ReadFile(file_url string) (content string, err error) {
	_, err = os.Stat(file_url)
	if err != nil {
		fmt.Println("cannot find the file ", file_url)
		return content, err
	}
	file, err := os.Open(file_url)
	defer file.Close()
	if err != nil {
		fmt.Println("cannot open the file ", file_url)
		return content, err
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	content = buf.String()

	return content, nil
}

func checkSemVer(version string) (ret bool, msg string) {
	ret = true
	str := strings.Split(version, ".")
	if len(str) != 3 {
		ret = false
	} else {
		for index := 0; index < len(str); index++ {
			i, err := strconv.Atoi(str[index])
			if err != nil {
				ret = false
				break
			} else {
				if i < 0 {
					ret = false
					break
				}
			}
		}
	}
	if ret == false {
		msg = fmt.Sprintf("%s is not a valid version format, please read 'SemVer v2.0.0'", version)
	}
	return ret, msg
}

func checkUnit(field reflect.Value, check string, parent string, err_msg []string) (bool, []string) {
	kind := field.Kind().String()
	switch kind {
	case "string":
		if check == "SemVer v2.0.0" {
			ok, msg := checkSemVer(field.String())
			if ok == false {
				err_msg = append(err_msg, fmt.Sprintf("%s : %s", parent, msg))
				return false, err_msg
			}
		}
		break
	default:
		break
	}
	return true, err_msg
}

func validateUnit(field reflect.Value, t_field reflect.StructField, parent string, err_msg []string) (bool, []string) {
	var mandatory bool
	if t_field.Tag.Get("mandatory") == "required" {
		mandatory = true
	} else {
		mandatory = false
	}

	kind := field.Kind().String()
	switch kind {
	case "string":
		if mandatory && (field.Len() == 0) {
			err_msg = append(err_msg, fmt.Sprintf("%s.%s is incomplete", parent, t_field.Name))
			return false, err_msg
		}
		break
	case "struct":
		if mandatory {
			return validateStruct(field, parent+"."+t_field.Name, err_msg)
		}
		break
	case "slice":
		if mandatory && (field.Len() == 0) {
			err_msg = append(err_msg, fmt.Sprintf("%s.%s is incomplete", parent, t_field.Name))
			return false, err_msg
		}
		valid := true
		for f_index := 0; f_index < field.Len(); f_index++ {
			if field.Index(f_index).Kind().String() == "struct" {
				var ok bool
				ok, err_msg = validateStruct(field.Index(f_index), parent+"."+t_field.Name, err_msg)
				if ok == false {
					valid = false
				}
			}
		}
		return valid, err_msg
		break
	case "int32":
		break
	default:
		break

	}

	check := t_field.Tag.Get("check")
	if len(check) > 0 {
		return checkUnit(field, check, parent+"."+t_field.Name, err_msg)
	}

	return true, err_msg
}

func validateStruct(value reflect.Value, parent string, err_msg []string) (bool, []string) {
	if value.Kind().String() != "struct" {
		fmt.Println("Program issue!")
		return true, err_msg
	}
	rtype := value.Type()
	valid := true
	for i := 0; i < value.NumField(); i++ {
		var ok bool
		field := value.Field(i)
		t_field := rtype.Field(i)
		ok, err_msg = validateUnit(field, t_field, parent, err_msg)
		if ok == false {
			valid = false
		}
	}
	if valid == false {
		err_msg = append(err_msg, fmt.Sprintf("%s is incomplete", parent))
	}
	return valid, err_msg
}
