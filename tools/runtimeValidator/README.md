## Runtime Validator       
      
The runtimeValidator aims to Verify if a runtime containers runs the bundle correctly, test its compliance to [opencontainers/specs](https://github.com/opencontainers/specs)      


### Summary for the impatient      
***Key note***           
Be sure to download [specs](htttps://github.com/opencontainers/specs) source code and install [runc](https://github.com/opencontainers/runc) first     

``` bash   
$    go get github.com/huawei-openlab/oct/tools/runtimeValidator                 #get source code       
$    cd $GOPATH/src/github.com/huawei-openlab/oct/tools/runtimeValidator         #change dir to spcsValidator
$    make                                                                        #build runtimeValidator      
$    ./runtimeValidator                                                          #run runtimeValidator     
```     
      

### Runtime Validator Quickstart
                
- **Using Tools**        

Tools used in runtimeValidator as plugins,
***Key Notes***   
[ocitools](github.com/zenlinTechnofreak/ocitools) are foked from [github.com/mrunalp](github.com/mrunalp/ocitools),         
adding some adaptor changes for oct.   

See [plugins/Makefile](./plugins/Makefile)     
       
``` Makefile    
all:    
  echo ${GOPATH}    
  echo "Installing plugin: github.com/zenlinTechnofreak/ocitools..."    
  set -e   
  go get github.com/zenlinTechnofreak/ocitools   
  go build github.com/zenlinTechnofreak/ocitools    
  go build github.com/zenlinTechnofreak/ocitools/cmd/runtimetest    
clean:    
  go clean    
  rm -rf ocitools runtimetest      
```    

- **Supportted Runtime**    
    
Only Support runc yet, going to support other runtimes in next step, changes should be existed in [factory](./factory)      


- **About Validation Cases**        

Cases are listed in [cases.conf](./cases.conf), It going to be rich, in the fomate of below: 
    
```   
process= --args=./runtimetest --args=vp --rootfs=rootfs --terminal=false;--args=./runtimetest --args=vp --rootfs=rootfs --terminal=false     
# result to generate two cases in [bundle](./bundle), should be bundle/process1 and bundle/process2   

```       

### What is good for runtimeValidator
1. Tools is used as plugins ,feel free to use any tools    
2. Cases can be free to be added into cases.conf
3. Uses goroutine, each go routine runs a case bundle to validate


### Next to Do 

1. Rich cases:
   Encrease the functionality of ocitools in [cmd/runtimetest](https://github.com/zenlinTechnofreak/ocitools/cmd/runtimetest)   
   Rich cases in [cases.conf](./cases.conf)    
2. Support other containers

### Reference
OCI specs on https://github.com/opencontainers/specs   

OCI runc on https://github.com/opencontainers/runc
