package main

import (
	"./specs"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

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

func main() {

	var configFile = flag.String("f", "", "input the config file, or default to config.json")
	flag.Parse()

	if len(*configFile) == 0 {
		*configFile = "./config.json"
	}

	var sp specs.Spec
	content, err := ReadFile(*configFile)
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
