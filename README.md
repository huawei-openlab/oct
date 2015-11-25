## OCI Test        
      
The ocitest aims to test if a runtime container is compliant with [opencontainers/specs](https://github.com/opencontainers/specs),     
It is a light weight testing framework, using ocitools and 3rd-party tools, managing configurable high coverage bundles as cases, supporting testing different runtimes.     


### Summary for the impatient      
***Key note***           
Be sure to download [specs](htttps://github.com/opencontainers/specs) source code and install [runc](https://github.com/opencontainers/runc) first     

``` bash   
$ go get github.com/huawei-openlab/oct                 #get source code       
$ cd $GOPATH/src/github.com/huawei-openlab/oct         #change dir to workspace 
$ make                                                 #build      
$ ./ocitest                                            #run     
```     
      

### OCI Test Quickstart       
       
- **Usage**      
       
``` sh      
$ ./ocitest --help
NAME:
   oci-testing - Utilities for OCI Testing,

    It is a light weight testing framework,
    using ocitools and 3rd-party tools, 
    managing configurable high coverage bundles as cases, 
    supporting testing different runtimes.

USAGE:
   ./ocitest [global options] command [command options] [arguments...]
   
VERSION:
   0.0.1
   
COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --runtime, -r "runc"         runtime to be tested, -r=runc or -r=rkt or -r=docker
   --help, -h			show help
   --generate-bash-completion	
   --version, -v		print the version
```    
      
- **Supportted runtime**     
      
|Name|Status|Testing Flow|
|------|----|------| ----- |
| runc | Supported| [Test bundles & runtime Validate](https://github.com/huawei-openlab/oct/blob/master/docs/static/runtime-validation-oci-standard.png) |
| rkt | Supported | [Test bundles converted by oci2aci & runtime Validate] (https://github.com/huawei-openlab/oct/blob/master/docs/static/runtime-validation-oci-standard2.png) |
| docker | Not currently being worked|[Test bundles converted by oci2docker & runtime Validate] (https://github.com/huawei-openlab/oct/blob/master/docs/static/runtime-validation-oci-standard2.png) |
      
- **Using Tools**        

Tools used by ocitest as plugins,
***Key Notes***        

[ocitools](github.com/zenlinTechnofreak/ocitools) are foked from [github.com/mrunalp](github.com/mrunalp/ocitools), adding some adaptor changes for oct.   

See [plugins/Makefile](./plugins/Makefile)     

- **About Test Cases**        

Cases are listed in [cases.conf](./cases.conf), as the fomate of bunldes, It is going to be rich, in the fomate of below: 
    
```   
process= --args=./runtimetest --args=vp --rootfs=rootfs --terminal=false;--args=./runtimetest --args=vp --rootfs=rootfs --terminal=false     
# result to generate two cases in [bundle](./bundle), should be bundle/process0 and bundle/process1,        
# and '--args=./runtimetest --args=vp --rootfs=rootfs --terminal=false' is params for ocitools generate   

```

### What is good for runtimeValidator       
1. Light weight testing freamwork      
2. High coverage test cases, configurable, easy to add cases
3. Tools is used as plugins ,feel free to use any 3rd-paty tools        
4. Uses goroutine, each go routine runs a case bundle to validate   
**Note**     
The ocitools are developed in [github.com/mrunalp](github.com/mrunalp/ocitools).  

### Next to Do 

1. Rich cases:        

   Encrease the functionality of ocitools in [cmd/runtimetest](https://github.com/zenlinTechnofreak/ocitools/tree/master/cmd/runtimetest)   
   Rich cases in [cases.conf](./cases.conf)    

2. Support other containers

### Reference
OCI specs on https://github.com/opencontainers/specs   

OCI runc on https://github.com/opencontainers/runc
