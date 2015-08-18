## Specstest-process-runc-
[Test Case](#testcase) Test runc when spec process

```
Owner: linzhinan@huawei.com
License: Apache 2.0
Group: Specstest/process/
```

The case has 1 host operation system(s):

'hostA' has 1 container(s) deployed.

The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
|hostA|ubuntu14.04|CPU 1, Memory 1GB, Disk 2G|specs(specstest)|"go build process_host.go ; ./process_host ; cd ./../../source/ ; ./test_process"|

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
|specstest|runc0.2||

After running the `Command` in each OS and container, we get two results.

* [process_out.json](#process_out) 


###TestCase
```
{
  "Name": "Specstest-process-runc",
  "Summary": "test based opencontainers/specs",
  "Owner": "linzhinan@huawei.com",
  "Description": "Test runc when spec process",
  "Group": "Specstest/process/",
  "License": "Apache 2.0",
  "Explains": "Test runc when spec process",
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
        "./source/process_host.go",
        "./source/process_guest.go",
        "./source/test_process.go"
      ],
      "Cmd": "go build process_host.go ; ./process_host ; cd ./../../source/ ; ./test_process",
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
        "/tmp/testtool/process_out.json"
      ]
    }
  ]
}
```

###process_out
```
{
  "Linuxspec.Spec.Process": {
    "terminal": {
      "true": "passed"
    },
    "user": {
      "uid": {
        "1": "passed"
      },
      "gid": {},
      "additionalGids": {
        "nil": "passed"
      }
    },
    "args": {
      "./process_guest": "passed"
    },
    "env": {
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin": "passed",
      "TERM=xterm": "passed"
    },
    "cwd": {
      "/testtool": "passed"
    }
  }
}
```


