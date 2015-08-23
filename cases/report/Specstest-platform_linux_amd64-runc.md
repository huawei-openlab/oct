## Specstest-platform_linux_amd64-runc-
[Test Case](#testcase) Test runc when linuxspec.Spec.Platform.OS = linux, linuxspec.Spec.Platform.Arch = amd64

```
Owner: linzhinan@huawei.com
License: Apache 2.0
Group: Specstest/platform_linux_amd64/
```

The case has 1 host operation system(s):

'hostA' has 1 container(s) deployed.

The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
|hostA|ubuntu14.04|CPU 1, Memory 1GB, Disk 2G|specs(specstest)|"go build platform_linux_amd64_host_deploy.go ; ./platform_linux_amd64_host_deploy ; cd ./../../source/ ;  ./test_platform_linux_amd64"|

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
|specstest|runc0.2| |

After running the `Command` in each OS and container, we get two results.

* [platform_linux_amd64_out.json](#platform_linux_amd64_out) 


###TestCase
```
{
  "Name": "Specstest-platform_linux_amd64-runc",
  "Summary": "test based opencontainers/specs",
  "Owner": "linzhinan@huawei.com",
  "Description": "Test runc when linuxspec.Spec.Platform.OS = linux, linuxspec.Spec.Platform.Arch = amd64",
  "Group": "Specstest/platform_linux_amd64/",
  "License": "Apache 2.0",
  "Explains": "Test runc when linuxspec.Spec.Platform.OS = linux, linuxspec.Spec.Platform.Arch = amd64 ",
  "Requires": [
    {
      "Class": "OperationOS",
      "Type": "os",
      "Distribution": "ubuntu",
      "Version": "14.04",
      "Resource": {
        "CPU": 1,
        "Memory": "1GB",
        "Disk": "2G"
      }
    },
    {
      "Class": "specstest",
      "Type": "container",
      "Distribution": "runc",
      "Version": "0.2"
    }
  ],
  "Deploys": [
    {
      "Object": "hostA",
      "Class": "OperationOS",
      "Files": [
        "./source/platform_linux_amd64_host_deploy.go",
        "./source/demo.go",
        "./source/test_platform_linux_amd64.go"
      ],
      "Cmd": "go build platform_linux_amd64_host_deploy.go ; ./platform_linux_amd64_host_deploy ; cd ./../../source/ ;  ./test_platform_linux_amd64",
      "Containers": [
        {
          "Object": "specs",
          "Class": "specstest"
        }
      ]
    }
  ],
  "Collects": [
    {
      "Object": "hostA",
      "Files": [
        "/tmp/testtool/platform_linux_amd64_out.json"
      ]
    }
  ]
}

```


###platform_linux_amd64_out
```
{
  "Linuxspec.Spec.Platform": {
    "OS": {
      "linux": "passed"
    },
    "Arch": {
      "amd64": "passed"
    }
  }
}
```


