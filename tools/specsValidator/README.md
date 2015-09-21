## The Specs test  aims to test OCP specs's compliance on runc

### Test Scope    

The specsValidator aims to test containers runtime compatible with  [opencontainers/specs](https://github.com/opencontainers/specs).Generates an OCI bundle that tests that the runtime is compliant
      


### Quickstart      
Just following steps below,      
``` bash
$    go get github.com/tools/godep                                               #install godep tool requested     
$    go get -d  -tags v0.1.1 github.com/huawei-openlab/oct/tools/specsValidator  #get source code           
$    make all BUILDTAGS=v0.1.1                                                   #build specsValidator      
$    ./specsValidator -rtags=apparmor                                            #run specsValidator     
$    cat report/linuxspec.json                                                   #get result       
# This steps will guiding to use the newest runc and specs, and test runc with buildtags=apparmor   
```     
      


### Summary for the impatient
      

- **Step 0: Prepare**

1. Install the go-lang env, set the GOPATH and GOROOT properly, reuquired go-lang version is V1.4.2      
2. Install godep tool, like the way below,
``` bash
$    go get github.com/tools/godep
```    
     
- **Step 1: Building**        

For using the specs commit Revxxx and runc commit Revyyy,      
There are two tags of specs can be used directly, v0.1.1 and predraft(commmit: 45ae53d4dba8e550942f7384914206103f6d2216)   
Not all of the specs version and runc version can match to run regularly, use the commit of runc and specs directly when the commmit of both of the runc and specs can be run reguarly together.

       
``` bash
$    go get -d -tags Revxxx github.com/huawei-openlab/oct/tools/specsValidator
$    cd $GOPATH/src/github.com/opencontainers/specs
$    git checkout Revxxx    
$    godep update github.com/opencontainers/specs
$    cd $GOPATH/src/github.com/huawei-openlab/oct/tools/specsValidator
$    make all BUILDTAGS=Revxxx    
```     
     

- **Step 2: Using**     
       
      
``` bash
$   su root
$   ./specsValidator -runc=<specified runc revision> -specs=<specified specs revision> -rtags=<specified runc build tags> -o=<output path>    
```      
For example,      
``` bash
$   ./specsValidator -specs=predraft -runc=v0.0.4 -rtags="apparmor"    
# or    
$   ./specsValidator -specs=v0.1.1 -runc=6b5a66f7e1444ac7776019a4bb8ad0b93584685d -rtags="apparmor"
```

Usage of specsValidator      
``` bash    
$   ./specsValidator --helpUsage of ./specsValidator:       
           
      -o="./report/": Specify filePath to install the test result linuxspec.json     
           
      -rtags="seccomp": Build tags for runc, should be one of seccomp/selinux/apparmor, keep empty to using seccomp      
            
      -runc="": Specify runc Revision from opencontainers/specs to be tested, in the form of commit id, keep empty to using the newest commit of [opencontainers/runc](https://github.com/opencontainers/runc       
           
      -specs="": Specify specs Revision from opencontainers/specs as the benchmark, in the form of commit id, keep empty to using the newest commit of [opencontainers/specs](https://github.com/opencontainers/specs)
```


    
- **Step 3: Getting Result**     

If specified the output path with "-o", get reuslt from the path specified,      
Else get result from [oct/tools/specsValidator/report](./report/), result is written in json format, total result is named linuxspec.json.
      
For example, in linuxspec.json
``` json
{
  "LinuxSpec.Linux.Namespaces": [
    {
      "testcasename": "TestMntPathEmpty",
      "input": {
        "type": "mount",
        "path": ""
      },
      "result": "passed"
    },
    {
      "testcasename": "TestMntPathUnempty",
      "input": {
        "type": "mount",
        "path": "/proc/1/ns/mnt"
      },
      "error": "time=\"2015-08-30T19:15:30+08:00\" level=warning msg=\"exit status 1\" \ntime=\"2015-08-30T19:15:30+08:00\" level=warning msg=\"open /sys/fs/cgroup/freezer/user/1000.user/c2.session/oct/freezer.state: no such file or directory\" \ntime=\"2015-08-30T19:15:30+08:00\" level=warning msg=\"open /sys/fs/cgroup/devices/user/1000.user/c2.session/oct/cgroup.procs: no such file or directory\" \ntime=\"2015-08-30T19:15:30+08:00\" level=fatal msg=\"Container start failed: [8] System error: invalid argument\" \nexit status 1",
      "result": "unspported"
    }
  ]
}
# Field "testcasename": name of testacse
# Field "input": input of the tetstcase to create config.json for containers, left value is the obj in [opencontainers/specs](https://github.com/opencontainers/specs), right value is the value of the obj in left.
# Field "error": output the err of the testcase
# Field "result": result of the testcase, it is value should be in {"passed", "failed", "unsupportd", "unknown"}
# "passed": testcase is passed
# "failed": testcase is failed
# "unsupported": the input of the config is not supported by runc yet
# "unkown": meet the unknown err, if anyone meet the result ,plz let me know
``` 

### Develop Progress

### Next to Do 

1. Rich cases for testing runtime of runc

### Reference
OCP specs on https://github.com/opencontainers/specs   

OCP runc on https://github.com/opencontainers/runc
