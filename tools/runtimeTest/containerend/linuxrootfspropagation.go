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
//

package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("/bin/bash", "-c", "mount -t tmpfs tmpfs /fspropagationtest")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux rootfs propagation test : mount fs in the container, %v", err)
	}
	cmd = exec.Command("/bin/bash", "-c", "touch  /fspropagationtest/fromcontainer.txt")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux rootfs propagation test : make new file in container error, %v", err)
	}
	cmd = exec.Command("/bin/bash", "-c", "ls  /fspropagationtest")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] linux rootfs propagation test : ls file in container error, %v", err)
	} else {
		fmt.Println(strings.TrimSpace(string(out)))
	}
}
