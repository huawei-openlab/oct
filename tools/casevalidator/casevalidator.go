package main

import (
	"flag"
	"fmt"
	"github.com/huawei-openlab/oct-engine/lib/libocit"
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

func main() {
	var caseDir = flag.String("d", "", "input the case dir")
	var caseFile = flag.String("f", "", "input the file url, case.tar.gz")
	var caseName = flag.String("n", "", "input the 'case name' in the case dir, if there were multiply cases in the case dir. You can use this with -d and -f.")
	var caseID = flag.String("id", "", "input the 'case id' provided by 'Test Case server', please make sure the the tcserver is running.")
	flag.Parse()

	var warning_msg []libocit.ValidatorMessage
	var err_msg []libocit.ValidatorMessage
	if len(*caseID) > 0 {
	} else if len(*caseFile) > 0 {
		libocit.ValidateByFile(*caseFile)
	} else if len(*caseDir) > 0 {
		warning_msg, err_msg = libocit.ValidateByDir(*caseDir, *caseName)
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
