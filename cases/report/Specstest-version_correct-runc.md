## Specstest-version_correct-runc-
[Test Case](#testcase) Test runc when spec version == pre-draft

```
Owner: linzhinan@huawei.com
License: Apache 2.0
Group: Specstest/version_correct/
```

The case has 1 host operation system(s):

'hostA' has 1 container(s) deployed.

The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
|hostA|ubuntu14.04|CPU 1, Memory 1GB, Disk 2G|specs(specstest)|"go build version_host.go ; ./version_host ; cd ./../source/ ; ./test_version_correct"|

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
|specstest|runc0.2| |

After running the `Command` in each OS and container, we get two results.

* [version_correct_out.json](#version_correct_out) 


###TestCase
```
{
  "Name": "Specstest-version_correct-runc",
  "Summary": "test based opencontainers/specs",
  "Owner": "linzhinan@huawei.com",
  "Description": "Test runc when spec version == pre-draft",
  "Group": "Specstest/version_correct/",
  "License": "Apache 2.0",
  "Explains": "Test runc when spec version == pre-draft",
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
      "Version": "0.2",
      "Files": [
        "./source/config.json"
      ]
    }
  ],
  "Deploys": [
    {
      "Object": "hostA",
      "Class": "OperationOS",
      "Files": [
        "./source/version_host.go",
        "./source/test_version_correct.go",
        "./source/demo.go"
      ],
      "Cmd": "go build version_host.go ; ./version_host ; cd ./../source/ ; ./test_version_correct",
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
        "/tmp/testtool/version_correct_out.json"
      ]
    }
  ]
}
```

###version_correct_out
```
{
  "Linuxspec.Spec.Version": {
    "pre-draft": "passed"
  }
}
```


