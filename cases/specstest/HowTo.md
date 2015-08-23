### How to Use Specstest   

Specstest aims to reach the reuslt as below
- Provide functionality to specified needed  [opencontainers/specs](https://github.com/opencontainers/specs)  version as a test benchmark.
- Do specs testing with the specified specs on different bundles/containers.
- Provide test reuslt of bundle/container with different specified specs.

The project is in developing and restructuring at present stage, and is not perfect yet,  so it have not supported feature of automatically specify specs version .      
 
 The project import two different version of  [opencontainers/specs](https://github.com/opencontainers/specs)
as local packages instead of pulling specified version of specs from  [opencontainers/specs](https://github.com/opencontainers/specs) automatically. Local packages are in [oct/cases/specstest/source/](./source/), as follows,
- [oct/cases/specstest/source/specs](./source/specs)      

   Specs version : commit 7414f4d3e90b5c22cae7c952d123e911c0cf94ba      

   All configuration of this version have been supported by runc. It can be use to test the support level with the this specs version.
- [oct/cases/specstest/source/specsnewest](./source/specsnewest)        

   Specs version : the newest commit in [opencontainers/specs](https://github.com/opencontainers/specs)


   It should be updated everyday from [opencontainers/specs](https://github.com/opencontainers/specs) by maintainers, runc can not work with it yet(Becasuse runc is also in progress).
            
* The project can only run each testcase by hand yet,  People who are interested in it can try it as follow steps *
- Do git clone the project to the GOPATH on your host machine,
``` 
git clone https://github.com/huawei-openlab/oct.git
```
- Install runc on your local machine    

Reference to [opencontainers/runc](https://github.com/opencontainers/runc)
-  Specify the specs version     

*Suggest to use the specs vesion: commit 7414f4d3e90b5c22cae7c952d123e911c0cf94ba*, just skip this step.       

To use newest specs version, do,
``` bash
mv oct/cases/specstest/source/specsnewest oct/cases/specstest/source/specs
mv oct/cases/specstest/source/config_newest.json oct/cases/specstest/source/config.json
```

####Prepare testcase   
The `official` way of doing this is mentioned here [HowTO](../../engine/tcserver/README.md).

But you can also do it by hand, taking the 'process' case for example:

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

You can also use the temporary script for convenience, taking the 'process' case for example:
   

``` bash
./startserver.sh
./runtestcase.sh -f process
```


- Get Result     

See the output of the last command in last step, it have pointed out the output place.
