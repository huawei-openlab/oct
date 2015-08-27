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
	"fmt"
	"os"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func validateLayout(path string) error {
        fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error accessing layout: %v", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("given path %q is not a directory", path)
	}
	var flist []string
	var cfgOK, runOK, rfsOK bool
	var config, runtime io.Reader
	walkLayout := func(fpath string, fi os.FileInfo, err error) error {
		rpath, err := filepath.Rel(path, fpath)
		if err != nil {
			return err
		}
		switch rpath {
		case ".":
		case ConfigFile:
			config, err = os.Open(fpath)
			if err != nil {
				return err
			}
			cfgOK = true
		case RuntimeFile:
			runtime, err = os.Open(fpath)
			if err != nil {
				return err
			}
			runOK = true
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
	return checkLayout(cfgOK, config, runOK, runtime, rfsOK, flist)
}

func checkLayout(cfgOK bool, config io.Reader, runOK bool, runtime io.Reader, rfsOK bool, files []string) error {
	defer func() {
		if rc, ok := config.(io.Closer); ok {
			rc.Close()
		}
		if rc, ok := runtime.(io.Closer); ok {
			rc.Close()
		}		
	}()
	if !cfgOK {
		return ErrNoConfig
	}
	if !runOK {
		return ErrNoRun 
	}
	if !rfsOK {
		return ErrNoRootFS
	}
	_, err := ioutil.ReadAll(config)
	if err != nil {
		return fmt.Errorf("error reading the layout: %v", err)
	}
	_, err = ioutil.ReadAll(runtime)
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
