The tools for open container test will be maintained here.

### [bundleValidator](bundleValidator/README.md) - specs/(bundle, config.json, runtime.json)
The `bundleValidator` is a `rolling release` tool which keeps updating with OCI specs.
It verifies if a `bundle` was a standard container, if a config.json/runtime.json is a standard configuration file.
(Current one is compliant with v0.1.0)

The option(genconfig/genruntime) in bundleValidator is used to generate a demo config.json/runtime.json.

### [specValidator](specsValidator/README.md) - specs/(runtime, config)
The `specValidator` verifies all the specs configurations, so it is a much bigger project
and we will only validate the `stable` spec release.
So far the tool is compliance with `45ae53d4dba8e550942f7384914206103f6d2216` commit.

## Make OCT better
We will also list other container testing tools here.
Please tell us and update this file by using the following format:

```
###[tool name](tool url) - tool function
tool description

```
