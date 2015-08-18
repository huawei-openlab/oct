## Specstest-linux-resources-memory-limited-
[Test Case](#testcase) Test runc  whether it supports the linux memory size limited

```
Owner: wangqilin2@huawei.com
License: Apache 2.0
Group: Specstest/linux_resources/
```

The case has 1 host operation system(s):

'hostA' has 1 container(s) deployed.

The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
|hostA|Ubuntu14.04|CPU 1, Memory 1GB, Disk 1G|specs(specstest)|"go build linux_resources_memory_limited_host.go ; ./linux_resources_memory_limited_host; go build linux_resources_memory_limited_runcrunner.go; cp ./linux_resources_memory_limited_runcrunner ./../../source/;"|

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
|specstest|runc V0.2| |

After running the `Command` in each OS and container, we get two results.

* [linux_resources_memory_limited.json](#linux_resources_memory_limited) 


###TestCase
```
{
  "Name": "Specstest-linux-resources-memory-limited",
  "Summary": "test linux_resources_memory_limited",
  "Owner": "wangqilin2@huawei.com",
  "Description": "Test runc  whether it supports the linux memory size limited",
  "Group": "Specstest/linux_resources/",
  "License": "Apache 2.0",
  "Explains": "Test runc by see whether it can constrain the  cgroup memory size",
  "Requires": [
    {
      "Class": "OperationOS",
      "Type": "os",
      "Distribution": "Ubuntu",
      "Version": "14.04",
      "Resource": {
        "CPU": 1,
        "Memory": "1GB",
        "Disk": "1G"
      }
    },
    {
      "Class": "specstest",
      "Type": "container",
      "Version": "runc V0.2"
    }
  ],
  "Deploys": [
    {
      "Object": "hostA",
      "Class": "OperationOS",
      "Cmd": "go build linux_resources_memory_limited_host.go ; ./linux_resources_memory_limited_host; go build linux_resources_memory_limited_runcrunner.go; cp ./linux_resources_memory_limited_runcrunner ./../../source/;",
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
      "Cmd": "cd  ./../../source/;./linux_resources_memory_limited_runcrunner"
    }
  ],
  "Collects": [
    {
      "Object": "hostA",
      "Files": [
        "/tmp/testtool/linux_resources_memory_limited.json"
      ]
    }
  ]
}```

###linux_resources_memory_limited
```
{"Linuxspec.Linux.Resources.Memory":{"Memory.Limit":"pass"}}
```


