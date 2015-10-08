The tools for open container test will be maintained here.

### [bundleValidator](bundleValidator/README.md) - specs/(bundle, config.json, runtime.json)
The `bundleValidator` is a `rolling release` tool which keeps updating with OCI specs.
It verifies if a `bundle` was a standard container, if a config.json/runtime.json is a standard configuration file.

The option(genconfig/genruntime) in bundleValidator is used to generate a demo config.json/runtime.json.

### [runtimeValidator](runtimeValidator/README.md) - specs/(runtime, config)
The `runtimeValidator` verifies if a runtime containers runs the bundle correctly, test its compliance to [opencontainers/specs](https://github.com/opencontainers/specs).

### [ociConvert](ociConvert/README.md) - convert from other container images
So far, image configurations compliants to AppC spec could be converted into OCI configurations.

## Make OCT better
We will also list other container testing tools here.
Please tell us and update this file by using the following format:

```
###[tool name](tool url) - tool function
tool description

```

## OCT with Other project
### OCT & OCT-Engine
It is very easy to [write](https://github.com/huawei-openlab/oct-engine/blob/master/cases/README.md) a testcase under the OCT-engine.
The testcase using bundleValidator as the testing tool is listed here:
https://github.com/huawei-openlab/oct-engine/blob/master/cases/bundle/config.json

