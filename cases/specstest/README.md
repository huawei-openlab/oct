## The Specs test  aims to test OCP specs's compliance on runc

### Test Scope    

The specstest aims to test containers runtime compatible with  [opencontainers/specs](https://github.com/opencontainers/specs).
It is test scope contains:   

- Validate the configuration of a bundle
- Generate and test a bundle/image cryptographic identity
- Generates an OCI bundle that tests that the runtime is compliant
- A tool for testing fetching/unpacking of OCI images

### How to Use Specs Test   

Specstest aims to reach the reuslt as below,
- Provide functionality to specified the needed  [opencontainers/specs](https://github.com/opencontainers/specs)  version to as a test benchmark.
- Do specs testing with the specified specs.
- Provide test reuslt of bundle/container with different specified specs.

The project is in developing and restructuring at present stage, it is not perfect yet, it have not  include automical version specified function.
In order to provide expericing specstest for user , the project just import two different version of  [opencontainers/specs](https://github.com/opencontainers/specs)
as a local package, int [oct/cases/specstest/source/](./source/) 
- [oct/cases/specstest/source/specs](./source/specs)      

   Specs version : commit 7414f4d3e90b5c22cae7c952d123e911c0cf94ba      

   All configuration of this version have been supported by runc. It can be use to test the support level with the this specs version.
- [oct/cases/specstest/source/specsnewest](./source/specsnewst)        

   Specs version : the newest commit in [opencontainers/specs](https://github.com/opencontainers/specs)


   It should be updated everyday from [opencontainers/specs](https://github.com/opencontainers/specs), runc can not work with it yet(Becasuse it is in progress also).
            
 ##### The project can only run each testcase by hand yet,  people who want experience it can use it as the way below,
- Do git clone the project to the GOPATH on your host machine,
``` go
git clone https://github.com/huawei-openlab/oct.git
```
- Install runc on your local machine    

Reference to [opencontainers/runc](https://github.com/opencontainers/runc)
-  Specify the specs version
Suggest to use the specs vesion: commit 7414f4d3e90b5c22cae7c952d123e911c0cf94ba, just skip this step.       

If u want to use newest specs version, your should do,
``` bash
mv oct/cases/specstest/source/specsnewest oct/cases/specstest/source/specs
mv oct/cases/specstest/source/config_newest.json oct/cases/specstest/source/config.json
```

- Prepare testcase
Testcase process for example,
``` bash
cd oct/cases/specstest/
mkdir tmp_dir
cd tmp_dir
cp -rf ../process .
cp -rf ../source .
tar czvf ../process.tar.gz *
cd ..
rm -rf tmp_dir
cp process.tar.gz oct/engine/scheduler
```

- Start test framework
``` bash
cd oct/engine
make
cd testsever
./testsever &
cd ../ocitd
./ocitd &
```
- Run testcase
``` bash
cd oct/engine/scheduler
./scheduler process.tar.gz
```
- Get Result
See the output fo the last command in last step, it have pointed out the output place.

### How to realize the specstest
It conclude the next steps,
- Uses the Specs in [opencontainers/specs](https://github.com/opencontainers/specs) as a benchmark,     
LinuxSpec is the the full specification for linux containers in the project.

- Convert  [config.json](./source/config.json) file to the obj LinuxSPec on the [opencontainers/specs](https://github.com/opencontainers/specs).       

Example:
```
var linuxspec *specs.LinuxSpec
linuxspec, err = configconvert.ConfigToLinuxSpec(filePath)
```

- Modify the test item of LinuxSpec obj to the testvalue in the host end.     

Example:
```
linuxspec.Spec.Root.Path = "rootfs_rootconfig"
linuxspec.Spec.Root.Readonly = true
```
- Convert LinuxSpec to the [config.json](./source/config.json) file.  

Example:
```
err = configconvert.LinuxSpecToConfig(filePath, linuxspec)
```

- Use rootfs specsed by the config.json and config.json to run runc.

- Run guest end test programme to check the result of the test item.

### Reference
OCP specs on https://github.com/opencontainers/specs   

OCP runc on https://github.com/opencontainers/runc
