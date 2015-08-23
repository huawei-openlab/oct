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
	CLASS        string
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
	TC     TestCase
	OSList []OS
}

type Deploy struct {
	Object     string
	Class      string
	Cmd        string
	Files      []string
	Containers []Container
	//if it was hostOS, the ID is the host OS ID
	ID string
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
	Version      string
	Resource     OSResource
	Files        []string
}

type Container struct {
	Object       string
	Class        string
	Cmd          string
	Files        []string
	Distribution string
	Version      string
}

type Collect struct {
	Object string
	Files  []string

	//if it was hostOS, the ID is the host OS ID
	ID string
}

type TestCase struct {
	Name        string
	Summary     string
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
	Data    string
}

//Set the object status, for example, set an HostA to 'running'
type TestingStatus struct {
	Object string
	Status string
}

type TestingCommand struct {
	//If it was hostOS, the ID is the hostOS ID
	ID     string
	Object string
	//Status: deploy, run
	Status  string
	Command string
}

//WHen filename is null, we just want to prepare a pure directory
func PreparePath(cachename string, filename string) (realurl string) {
	var dir string
	if filename == "" {
		dir = cachename
	} else {
		realurl = path.Join(cachename, filename)
		dir = path.Dir(realurl)
	}
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

func SendFile(post_url string, file_url string, params map[string]string) (ret HttpRet) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	filename := path.Base(file_url)
	//'tcfile': testcase file
	fileWriter, err := bodyWriter.CreateFormFile("tcfile", filename)
	if err != nil {
		ret.Status = "Failed"
		ret.Message = "error writing to buffer"
		return ret
	}
	_, err = os.Stat(file_url)
	if err != nil {
		ret.Status = "Failed"
		ret.Message = "error in stat file " + file_url
		return ret
	}

	fh, err := os.Open(file_url)
	if err != nil {
		ret.Status = "Failed"
		ret.Message = "error in open file " + file_url
		return ret
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		ret.Status = "Failed"
		ret.Message = "error in copy file " + file_url
		return ret
	}

	for key, val := range params {
		fmt.Println("key  ", key, "  val  ", val)
		_ = bodyWriter.WriteField(key, val)
	}
	//	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()
	request, err := http.NewRequest("POST", post_url, bodyBuf)
	if err != nil {
		ret.Status = "Failed"
		ret.Message = "error in get new request"
		return ret
	}
	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		ret.Status = "Failed"
		ret.Message = "error in send new request"
		return ret
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ret.Status = "Failed"
		ret.Message = "error in reading response"
	} else {
		ret.Status = "OK"
		ret.Message = string(resp_body)
	}
	return ret

}

func SendCommand(apiurl string, b []byte) (ret HttpRet) {
	body := bytes.NewBuffer(b)
	resp, perr := http.Post(apiurl, "application/json;charset=utf-8", body)
	if perr != nil {
		fmt.Println(perr)
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
		resp.Body.Close()
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

func ReceiveFile(w http.ResponseWriter, r *http.Request, cache_url string) (real_url string, params map[string]string) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("tcfile")

	params = make(map[string]string)

	if r.MultipartForm != nil {
		for key, val := range r.MultipartForm.Value {
			//Must use val[0]
			params[key] = val[0]
		}
	}

	if err != nil {
		fmt.Println("Cannot find the tc file")
		return real_url, params
	}
	defer file.Close()

	real_url = PreparePath(cache_url, handler.Filename)
	f, err := os.Create(real_url)
	if err != nil {
		fmt.Println("Cannot create the file ", real_url)
		//TODO: better system error
		http.Error(w, err.Error(), 500)
		return real_url, params
	}
	defer f.Close()
	io.Copy(f, file)

	return real_url, params
}

// file name filelist is like this: './source/file'
func TarFileList(filelist []string, case_dir string, object_name string) (tar_url string) {
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
	tar_url = TarFileList(files, case_dir, case_name)
	return tar_url
}

func UntarFile(cache_url string, filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		fmt.Println("cannot find the file ", filename)
		return
	}

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
		fw, err := os.OpenFile(target_url, os.O_CREATE|os.O_WRONLY, os.FileMode(h.Mode))
		if err != nil {
			//Dir for example
			continue
		} else {
			io.Copy(fw, tr)
			fw.Close()
		}
	}
}

func ReadCaseFromTar(tar_url string) (content string) {
	_, err := os.Stat(tar_url)
	if err != nil {
		fmt.Println("cannot find the file ", tar_url)
		return content
	}

	fr, err := os.Open(tar_url)
	if err != nil {
		fmt.Println("fail in open file ", tar_url)
		return content
	}
	defer fr.Close()
	gr, err := gzip.NewReader(fr)
	if err != nil {
		fmt.Println("fail in using gzip")
		return content
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
		fileSuffix := path.Ext(h.Name)
		if fileSuffix == ".json" {
			var tc TestCase
			buf := bytes.NewBufferString("")
			buf.ReadFrom(tr)
			file_content := buf.String()
			json.Unmarshal([]byte(file_content), &tc)
			if len(tc.Name) > 1 {
				content = file_content
				break
			} else {
				continue
			}
		}
	}

	return content
}

//file_url is the default file, suffix is the potential file
func ReadTar(tar_url string, file_url string, suffix string) (content string) {
	_, err := os.Stat(tar_url)
	if err != nil {
		fmt.Println("cannot find the file ", tar_url)
		return content
	}

	fr, err := os.Open(tar_url)
	if err != nil {
		fmt.Println("fail in open file ", tar_url)
		return content
	}
	defer fr.Close()
	gr, err := gzip.NewReader(fr)
	if err != nil {
		fmt.Println("fail in using gzip")
		return content
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
		if len(suffix) > 0 {
			fileSuffix := path.Ext(h.Name)
			if fileSuffix == suffix {
				buf := bytes.NewBufferString("")
				buf.ReadFrom(tr)
				content = buf.String()
				break
			}
		}
		if len(file_url) > 0 {
			if h.Name == file_url {
				buf := bytes.NewBufferString("")
				buf.ReadFrom(tr)
				content = buf.String()
				break
			}
		}
	}

	return content
}

func MD5(data string) (val string) {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))

}
