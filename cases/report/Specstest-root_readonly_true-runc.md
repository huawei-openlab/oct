## Specstest-root_readonly_true-runc-
[Test Case](#testcase) Test runc when spec root:readonly == true

```
Owner: linzhinan@huawei.com
License: Apache 2.0
Group: Specstest/root_readonly_true/
```

The case has 1 host operation system(s):

'hostA' has 1 container(s) deployed.

The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
|hostA|ubuntu14.04|CPU 1, Memory 1GB, Disk 2G|specs(specstest)|"go build root_readonly_true_host.go ; ./root_readonly_true_host"|

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
|specstest|runc0.2| |

After running the `Command` in each OS and container, we get two results.

* [readonly_true_out.json](#readonly_true_out) 


###TestCase
```
{
  "Name": "Specstest-root_readonly_true-runc",
  "Summary": "test based opencontainers/specs",
  "Owner": "linzhinan@huawei.com",
  "Description": "Test runc when spec root:readonly == true",
  "Group": "Specstest/root_readonly_true/",
  "License": "Apache 2.0",
  "Explains": "Test runc when spec root:readonly == true",
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
        "./source/root_readonly_true_host.go",
        "./source/root_readonly_true_guest.go"
      ],
      "Cmd": "go build root_readonly_true_host.go ; ./root_readonly_true_host",
      "Containers": [
        {
          "Object": "specs",
          "Class": "specstest"
        }
      ]
    }
  ],
  "Run": [
    {
      "Object": "hostA",
      "Class": "OperationOS",
      "Cmd": "cd ./../../source/ ; runc"
    }
  ],
  "Collects": [
    {
      "Object": "hostA",
      "Files": [
        "/tmp/testtool/readonly_true_out.json"
      ]
    }
  ]
}
```

###readonly_true_out
```
{
  "Linuxspec.Spec.Root.Readonly": {
    "true": "passed"
  }
}
```


