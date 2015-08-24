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
)

func main() {
	cmd := exec.Command("/bin/sh", "-c", "echo \"FilesystemSupportValidationContent\" > /mountTest/test.txt")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] Mount test create file inside container failed , err = %v...", err)
	}
	cmd = exec.Command("/bin/sh", "-c", "cat  /mountTest/test.txt")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("[Specstest] Mount test read file inside container failed , err = %v...", err)
	}
	fmt.Println(string(output))
}
