## Specstest-mount_fstype_tmpfs-
[Test Case](#testcase) Test runc whether it can mount the tmpfs filesystem

```
Owner: wangqilin2@huawei.com
License: Apache 2.0
Group: Specstest/mount_fstype/
```

The case has 1 host operation system(s):

'hostA' has 1 container(s) deployed.

The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
|hostA|ubuntu14.04|CPU 1, Memory 1GB, Disk 2G|specs(specstest)|"go build mount_fstype_host.go;./mount_fstype_host --fs tmpfs ; cp mount_fstype_runcrunner.sh ./../../source/ "|

The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
|specstest|runc0.2||

After running the `Command` in each OS and container, we get two results.

* [mount_fstype_out.json](#mount_fstype_out) 


###TestCase
```
{
  "Name": "Specstest-mount_fstype_tmpfs",
  "Summary": "test based opencontainers/specs",
  "Owner": "wangqilin2@huawei.com",
  "Description": "Test runc whether it can mount the tmpfs filesystem",
  "Group": "Specstest/mount_fstype/",
  "License": "Apache 2.0",
  "Explains": "Test whether the container supports to mount a filesystem",
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
      "Cmd": "go build mount_fstype_host.go;./mount_fstype_host --fs tmpfs ; cp mount_fstype_runcrunner.sh ./../../source/ ",
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
      "Cmd": " cd  ./../../source/;./mount_fstype_runcrunner.sh tmpfs"
    }
  ],
  "Collects": [
    {
      "Object": "hostA",
      "Files": [
        "/tmp/testtool/mount_fstype_out.json"
      ]
    }
  ]
}
```

###mount_fstype_out
```
{
  "Linuxspec.Mounts.Filesystem": {
    "tmpfs": "success"
  }
}
```


