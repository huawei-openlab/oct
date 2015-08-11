package libocit

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type OS struct {
	ID           string
	Distribution string
	Version      string
	Arch         string
	CPU          int64
	Memory       int64
	IP           string
	Status       string
	Object       string
}

type Task struct {
	ID     string
	OSList []OS
}

//TODO add a 'casename' ?
type Deploy struct {
	Object     string
	Class      string
	Cmd        string
	Files      []string
	Containers []Container

	ResourceID string
}

//FIXME: the type is not consistent
type OSResource struct {
	CPU    int
	Memory string
	Disk   string
}

type Require struct {
	Class        string
	Type         string
	Distribution string
	Version      int
	Resource     OSResource
	Files        []string
}

type Container struct {
	Object       string
	Class        string
	Cmd          string
	Files        []string
	Distribution string
	Version      int
}

type Collect struct {
	Object string
	Files  []string

	ResourceID string
}

type Resource struct {
	//TODO: put following to a struct and make a hash?
	ID     string //returned
	Status bool   //whether it is available
	Msg    string //return value from server

	Req  Require
	Used bool
}

type TestCase struct {
	Name        string
	Version     string
	License     string
	Group       string
	Owner       string
	Description string
	Sources     []string
	Requires    []Require
	Deploys     []Deploy
	Run         []Deploy
	Collects    []Collect
}

type HttpRet struct {
	Status  string
	Message string
}

type TestingStatus struct {
	Object string
	Status string
}

type TestingCommand struct {
	ID     string
	Object string
	//Command: deploy, run
	Command string
}

func PreparePath(cachename string, filename string) (realurl string) {
	realurl = path.Join(cachename, filename)
	dir := path.Dir(realurl)
	p, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(dir, 0777)
		}
	} else {
		if p.IsDir() {
			return realurl
		} else {
			os.Remove(dir)
			os.MkdirAll(dir, 0777)
		}
	}
	return realurl
}

func SendFile(post_url string, file_url string, filename string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//'tcfile': testcase file
	fileWriter, err := bodyWriter.CreateFormFile("tcfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return
	}
	_, err = os.Stat(file_url)
	if err != nil {
		fmt.Println("error in stat file")
		return
	}

	fh, err := os.Open(file_url)
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println("error copy file")
		return
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(post_url, contentType, bodyBuf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(resp.Status)
	fmt.Println(resp_body)
}

func SendCommand(apiurl string, b []byte) (ret HttpRet) {
	body := bytes.NewBuffer(b)
	resp, perr := http.Post(apiurl, "application/json;charset=utf-8", body)
	defer resp.Body.Close()
	if perr != nil {
		ret.Status = "Failed"
		ret.Message = "err in posting"
	} else {
		result, berr := ioutil.ReadAll(resp.Body)
		if berr != nil {
			ret.Status = "Failed"
			ret.Message = "err in reading the response of the posting"
		} else {
			json.Unmarshal([]byte(result), &ret)
		}
	}
	return ret
}

//TODO: add err para?
func ReadFile(file_url string) (content string) {
	_, err := os.Stat(file_url)
	if err != nil {
		fmt.Println("cannot find the file ", file_url)
		return content
	}
	file, err := os.Open(file_url)
	defer file.Close()
	if err != nil {
		fmt.Println("cannot open the file ", file_url)
		return content
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	content = buf.String()

	return content
}

func ReceiveFile(w http.ResponseWriter, r *http.Request, cache_url string) (real_url string, handle_name string) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("tcfile")
	fmt.Println(handler.Filename)
	if err != nil {
		fmt.Println("Cannot find the tc file")
		return real_url, handle_name
	}
	defer file.Close()

	real_url = PreparePath(cache_url, handler.Filename)
	f, err := os.Create(real_url)
	if err != nil {
		fmt.Println("Cannot create the file ", real_url)
		//TODO: better system error
		http.Error(w, err.Error(), 500)
		return real_url, handle_name
	}
	defer f.Close()
	io.Copy(f, file)

	handle_name = handler.Filename
	return real_url, handle_name
}

// file name filelist is like this: './source/file'
func TarFilelist(filelist []string, case_dir string, object_name string) (tar_url string) {
	tar_url = path.Join(case_dir, object_name) + ".tar.gz"
	fw, err := os.Create(tar_url)
	if err != nil {
		fmt.Println("Failed in create tar file ", err)
		return tar_url
	}
	defer fw.Close()
	gw := gzip.NewWriter(fw)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for index := 0; index < len(filelist); index++ {
		source_file := filelist[index]
		fi, err := os.Stat(path.Join(case_dir, source_file))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fr, err := os.Open(path.Join(case_dir, source_file))
		if err != nil {
			fmt.Println(err)
			continue
		}
		h := new(tar.Header)
		h.Name = source_file
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()
		err = tw.WriteHeader(h)
		_, err = io.Copy(tw, fr)
	}
	return tar_url
}

func GetDirFiles(base_dir string, dir string) (files []string) {
	files_info, _ := ioutil.ReadDir(path.Join(base_dir, dir))
	for _, file := range files_info {
		if file.IsDir() {
			sub_files := GetDirFiles(base_dir, path.Join(dir, file.Name()))
			for _, sub_file := range sub_files {
				files = append(files, sub_file)
			}
		} else {
			files = append(files, path.Join(dir, file.Name()))
		}
	}
	return files

}

func TarDir(case_dir string) (tar_url string) {
	files := GetDirFiles(case_dir, "")
	case_name := path.Base(case_dir)
	tar_url = TarFilelist(files, case_dir, case_name)
	return tar_url
}

func UntarFile(cache_url string, filename string) {
	fr, err := os.Open(filename)
	if err != nil {
		fmt.Println("fail in open file ", filename)
		return
	}
	defer fr.Close()
	gr, err := gzip.NewReader(fr)
	if err != nil {
		fmt.Println("fail in using gzip")
		return
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		target_url := PreparePath(cache_url, h.Name)
		fw, _ := os.OpenFile(target_url, os.O_CREATE|os.O_WRONLY, os.FileMode(h.Mode))
		defer fw.Close()

		io.Copy(fw, tr)
	}
}

func MD5(data string) (val string) {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))

}
