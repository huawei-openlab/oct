## The Specs test  aims to test OCP specs's compliance on runc

### Test Scope    

The specsValidator aims to test containers runtime compatible with  [opencontainers/specs](https://github.com/opencontainers/specs).Generates an OCI bundle that tests that the runtime is compliant

### How to Use

- Prepare

1. Install the go-lang env, set the GOPATH and GOROOT properly, reuquired go-lang version is V1.4.2      
2. Install godep tool, like the way below,
``` bash
$       go get github.com/tools/godep
```
- Building        
It just guide people to use the specsValidator tool directly, but not using it in the whole project,         
if anyone wan to use it accross the whole proeject, please go through the [oct/README.md](./../../README.md)
``` bash
$       go get -d  github.com/huawei-openlab/oct/tools/specsValidator
$       cd $GOPATH/github.com/huawei-openlab/oct/tools/specsValidator
$       godep go build .
```
Binary "specs" is buit now.
- Using    
Use binary "specs" to run,
Usage of ./specs:        
  -o="./report/": Specify filePath to install the test result linuxspec.json     
  -runc="": Specify runc Revision from opencontainers/specs to be tested, in the form of commit id        
  -specs="": Specify specs Revision from opencontainers/specs as the benchmark, in the form of commit id

``` bash
$       su root
$       ./specsValidator -runc=<specified runc revision> -specs=<specified specs revision> -o=<output path>
```
#### Specs version should be older than 45ae53d4dba8e550942f7384914206103f6d2216, becasue of the huge change of adding runtime.config in specs.

- Getting Result     
Get result from [oct/tools/specsValidator/report](./report/), result is written in json format.
For example, in namespace_out.json
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

It can support to resive the specified specs revision and runc revision now, But can not support to use the newest specs version ,becasuse the specs have the huge change in adding the runtime.json.

### Next to Do 

1. Support the newest version of specs.
2. Provide overall testcase of the specs.
2. Support other containers.

### Reference
OCP specs on https://github.com/opencontainers/specs   

OCP runc on https://github.com/opencontainers/runc
