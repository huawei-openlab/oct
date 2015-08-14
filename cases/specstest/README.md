## The Specs test  aims to test OCP specs's compliance on runc

### Test Scope    

The specstest aims to test containers runtime compatible with  [opencontainers/specs](https://github.com/opencontainers/specs).
It is test scope contains:   

- Validate the configuration of a bundle
- Generate and test a bundle/image cryptographic identity
- Generates an OCI bundle that tests that the runtime is compliant
- A tool for testing fetching/unpacking of OCI images

### About Json Files

The JSON file should be called by the test scheduler to do test jobs, include serveral parts:
``` json
{"Description": {}, "Requires": {}, "Deploy":{}, "Run": {}, "Collect": {}}
```
For more information, refers to [oct/README.md](./../../README.md).

### Specstest Steps
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
