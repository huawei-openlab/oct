// Copyright 2014 Google Inc. All Rights Reserved.
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
       "log"
       "os/exec"
)

func init() {

}

func iperf_Server() {
	//exec shell in host machine to get docker container id
	cmd := exec.Command("/bin/sh", "-c", "docker ps -q")
	short_id, err := cmd.Output()
	if err != nil {
		log.Fatalf("Get container short id error %v", err)
		return
	}

	cmd = exec.Command("/bin/sh", "-c", "docker inspect -f '{{.id}}' "+string(short_id))
	full_id, err := cmd.Output()
	if err != nil {
	    log.Fatalf("Get container full id error %v", err)
	    return
	}
	fmt.Println("The container id of iperf server docker is %v", string(full_id))
	
	//Container name : "/docker/container id"
	containerName := "/docker/" + string(full_id)
	
	exec.Command("/bin/sh", "-c", "docker attach %v", string(full_id))
	
	perfServer, err := exec.Command("/bin/sh", "-c", "docker inspect -f '{{.NetworkSettings.IPAddress}}' ")
	
	//The output is JSON format.
	exec.Command("/bin/sh", "-c", "iperf -s -j")

}

func main() {
    iperf_Server()	
}
