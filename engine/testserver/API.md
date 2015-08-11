# Test Server 
The Test Server is used to manage the server resource.

The [configuration](#configs "configuration") file is used for the OCT users to set his/her own configuration.

The [openAPIs](#apis "APIs") are used for the developers to integrate the 'Test Server' with other services.

The [attributes](#attributes "attributes") are listed at the end of this document.

##Configs
|Key|Type|Description|Example|
|------|----|------| ----- |
| Port | int | The port of the Test Server.| 8001 |
| ServerListFile | string | For simple cases, we put all the server infos into a file.| ["servers.conf"](#servers) |
| CacheDir | string | The cache dir where the temporal files stored.| "/tmp/testserver_cache" |
| Debug | bool | Print the debug information on the screen| true, default to false |

```
{
	"Port": 8001,
	"ServerListFile": "servers.conf",
	"CacheDir" : "/tmp/testserver_cache", 
	"Debug": true
}
```
###servers

```
{
"Servers":
[
	{
	"Distribution": "Ubuntu",
	"Version": "12.04",
        "Arch": "x86_64",
        "ID": "0001",
        "CPU": 2,
        "Memory": 4,
        "IP": "192.168.0.1"
	},
	{
        "Distribution": "CentOS",
        "Version": "7",
        "Arch": "x86_64",
        "ID": "0002",
        "CPU": 1,
        "Memory": 3,
        "IP": "192.168.0.2"
	}
 ]
}
```

##APIs

|Method|Path|Summary|Description|
|------|----|------|-----------|
| POST | `/task` | [Upload files](#upload "Upload") | Upload the test case files, name: taskID.tar.gz. Used by 'Scheduler'|
| GET | `/:ID/result` | [Get test result](#result "Result") | Fetch the case result file. %taskid-result.tar.gz. Used by 'Schedular'|
| GET | `/:ID/status` | [Get testing status result](#get-status "Get Status") | Fetch the testing status. Used by 'Schedular'|
| POST | `/:ID/status` | [Set the testing status result](#set-status "Set Status") | Set the testing status. Used by 'OCTD'|
| GET | `/os` | [Resource](#resource "Resource") | Get the host resource on the server.|
| POST| `/os` | [Add resource](#add "Add resource") | Add a new host OS node, usually done by OCTD automaticly.|
| GET | `/os/:OSID` | [Detailed Resource](#details "Details") | Get the detailed information of a host OS.|

###upload

Upload the test files, name: taskID.tar.gz
```
　　Content-Disposition: form-data; name="tcfile"; filename="%taskid.tar.gz"
　　Content-Type: application/x-gzip
```

**Response**

```
  { "Status": "ok",
    "Message": "success in receiving task files"
  }

```

###result

```
GET /:ID/result
```

**Response**

The result files, listed in the test case 'Collect' session, for example:

```
  {Collect: ["./source/a.json", "./source/b.json"]}
```

The returned value should be like this:

```
  taskID.tar.gz 
         |____ source
                 |________a.json
                 |________b.json
                 
```

###get-status

```
GET /:ID/status
```

**Response**

"Applying/Allocated/Deploying/Deployed/Running/Finish"

``` 
  {"ID" "0001",
   [
    {
//OS ID
     "ID" : "10001",
     "Object": "HostA",
     "Distribution": "CentOS",
      ....
     "Status": "Deploying"
    },
    {
     "ID" : "10002",
     "Object": "HostB",
     "Distribution": "CentOS",
      ....
     "Status": "Deployed"
    }
   ]
  }
```

###set-status

```
POST /:ID/status
```

**Input**

| *Name* | *Type* | *Description* |
| -------| ------ | --------- |
| Status | string | The status of the task.|
| Object | string | The object name in the task.|

**Example**

```
  curl -i -d '{"Status":"deploy", "Object": "hostA"}'  /10002/status
```

**Response**


``` 
  {"Status": "OK",
   "Message": "The status is changed"
  }
```

####resource

```
GET /os
```
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
    "Version": "7",
    "CPU": 16,
    "Memory": 100000,
    "IP": 192.168.100.1,
    "Status": "free",
  },
  
  {
    "ID": "1002",
    "Distribution": "CentOS",
    "Version": "7",
    "CPU": 32,
    "Memory": 100000,
    "Disk": 10000,
    "IP": "192.168.100.2",
    "Status": "locked",
  }
]
```


###add
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
  curl -i -d '{"Distribution":"CentOS", "Version": "7",
              "Arch": "X86_64", "CPU": 32, "Memory": 100000, 
              "Disk": 10000, "IP": "192.168.100.2"]}'  /os
```

###details

```
GET /os/:OSID
```

On the TODO list, do not expose the details for now.


##Attributes

**ID**
`The task ID of the tasks.`

**OSID**
`The os ID in the server`

**CaseName**
`The name of the test case, same with the API URL for fetching the event.`

**Status**
`'OK/Failed'`

**Message**
`The message, especially the error message from the test server.`
