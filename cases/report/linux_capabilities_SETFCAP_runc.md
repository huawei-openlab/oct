## Specstest-linux-capabilities-SETFCAP-
[Test Case](#testcase) Test runc  whether it supports the linux capabilities SETFCAP

```
Owner: wangqilin2@huawei.com
License: Apache 2.0
Group: Specstest/linux_capabilities/
```

The case has 1 host operation system(s):

'hostA' has 1 container(s) deployed.

The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
|hostA|ubuntu14.04|CPU 1, Memory 1GB, Disk 2G|specs(specstest)|"go build linux_capabilities_SETFCAP_host.go ; ./linux_capabilities_SETFCAP_host;"|

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
|specstest|runc0.2| |

After running the `Command` in each OS and container, we get two results.

* [linux_capabilities_SETFCAP.json](#linux_capabilities_SETFCAP) 


###TestCase
```
{
  "Name": "Specstest-linux-capabilities-SETFCAP",
  "Summary": "test linux-capabilities-SETFCAP",
  "Owner": "wangqilin2@huawei.com",
  "Description": "Test runc  whether it supports the linux capabilities SETFCAP",
  "Group": "Specstest/linux_capabilities/",
  "License": "Apache 2.0",
  "Explains": "Test runc by set the capabilities SETFCAP to an other file",
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
      "Cmd": "go build linux_capabilities_SETFCAP_host.go ; ./linux_capabilities_SETFCAP_host;",
       "Files": [
        "./source/linux_capabilities_SETFCAP_guest.go",
        "./source/linux_capabilities_SETFCAP_host.go"
      ],
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
      "Cmd": " cd  ./../../source/;runc"
    }
  ],
  "Collects": [
    {
      "Object": "hostA",
      "Files": [
        "/tmp/testtool/linux_capabilities_SETFCAP.json"
      ]
    }
  ]
}
```

###linux_capabilities_SETFCAP
```
{
  "Linuxspec.Linux.Capabilities": {
    "SETFCAP": "pass"
  }
}
```


