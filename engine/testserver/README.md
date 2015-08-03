# Test Server 
The Test Server is used to manage the server resource.
Receive the 'test request' from 'Schedular', deploy the case to the avaliable host OS and run the test.
The Test Server provides some Restful APIs to the this framework user/developer.
By using the APIs, the server admin could also monitor the running status of the whole server nodes.

## APIs

* [Query the server resource](#query "query") 
* [Add a new host OS node to the server](#add "Add")


###Query

```
GET /os
```
List all the available OS for the specified parameters.

**Parameters**

| *Name* | *Type* | *Description* |
| -------| ------ | --------- |
| Distribution |	string | The distribution name. For example, CentOS, openSUSE, Ubuntu. Default: all |
| Version | string | The distribution version. For example, '7' (CentOS7), Default:all |
| Arch | string | The architecture of the OS. For example, x86_64, arm64. Default:all |
| CPU | int | The minimal number of the cpu. Default: 0 |
| Memory | int | The minimal amount of the memory (MB). Default: 0 |
| Disk | int | The minimal amount of the disk space (GB). Default : 0 |
| Page | int | The common parameter. Default: 0 |
| PageSize | int | The page size. Default to 10 |

**Response**

```
[
  {
    "ID": "1000",
    "Distribution": "CentOS",
    "Version": 7,
    "CPU": 16,
    "Memory": 100000,
    "IP": 192.168.100.1,
    "Status": "free",
  },
  {
    "ID": "1002",
    "Distribution": "CentOS",
    "Version": 7,
    "CPU": 32,
    "Memory": 100000,
    "Disk": 10000,
    "IP": "192.168.100.2",
    "Status": "locked",
  }
]
```

**Status**
```
  free : avaiable
  locked: in use
```

###Add
```
POST /os
```
Add a new host OS node to the Test Server. Most time it is done automaticly when an 'ocitd' daemon start running in a new node.

**Input**

| *Name* | *Type* | *Description* |
| -------| ------ | --------- |
| Distribution | string | The distribution name, mandatory |
| Version | string | The distribution version, mandatory |
| Arch | string | The architecture, mandatory |
| CPU | int | The minimal number of the cpu, mandatory |
| Memory | int | The amount of the memory (MB), mandatory |
| Disk | int | The amount of the disk space (GB), mandatory |
| IP | string | The IP address. (mandatory |

**Example**
```
  curl -i -d '{"Distribution":"CentOS", "Version": 7,
              "Arch": "X86_64", "CPU": 32, "Memory": 100000, 
              "Disk": 10000, "IP": "192.168.100.2"]}'  /os
```
