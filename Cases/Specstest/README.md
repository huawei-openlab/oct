## The Specs test  aims to test OCP specs's compliance on runc
The testcases should be written in JSON format, Each testcase should conclude a souce fold and a JSON file
### JSON file
The JSON file should be called by the test scheduler to do test jobs, include serveral parts:
```
{"Description": {}, "Requires": {}, "Deploy":{}, "Run": {}, "Collect": {}}
```
For more information, refers to [ocp-testing/README.md](./../../README.md).

### How to Test Specs
It conclude the next steps,
- Uses the Specs in [opencontainers/specs](https://github.com/opencontainers/specs) as a benchmark,     
LinuxSpec is the the full specification for linux containers in the project.

- Convert  [config.json](./source/config.json) file to the obj LinuxSPec on the [opencontainers/specs](https://github.com/opencontainers/specs).       

example:
```
var linuxspec *specs.LinuxSpec
linuxspec, err = configconvert.ConfigToLinuxSpec(filePath)
```

- Modify the test item of LinuxSpec obj to the testvalue in the host end.     

example:
```
linuxspec.Spec.Root.Path = "rootfs_rootconfig"
linuxspec.Spec.Root.Readonly = true
```
- Convert LinuxSpec to the [config.json](./source/config.json) file.  

example:
```
err = configconvert.LinuxSpecToConfig(filePath, linuxspec)
```

- Use rootfs specsed by the config.json and config.json to run runc.

- Run guest end test programme to check the result of the test item.

### Reference
OCP specs on https://github.com/opencontainers/specs   

OCP runc on https://github.com/opencontainers/runc
