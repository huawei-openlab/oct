#OCTD 
The 'OCID' is the program installed on the hostOS, used to communicate with the 'Test Server'.
(TODO:  All 'OCITD' should rename by OCTD since the project name changed)

- The [configuration](#configs "configuration") file is used for the OCT users to set his/her own configuration.
- The [openAPIs](#apis "APIs") are used for the 'OCT' developer to integrate 'OCTD' with the 'Test Server'.
- The [attributes](#attributes "attributes") are listed at the end of this document.

##Configs
|Key|Type|Description|Example|
|------|----|------| ----- |
| TSurl | string | The url of the Test Server, with the port| "http://localhost:8001" |
| Port | int | Port of the OCTD daemon| 9001 |
| CacheDir | string | Where are all the temporary files stored| "/tmp/octd-cache" |
| Debug | bool | Print the debug information on the screen| true, default to false |

```
{
	"TSurl": "http://localhost:8001",
	"Port": 9001,
	"CacheDir": "/tmp/octd-cache",
	"Debug": false
}
```

##APIs

|Method|Path|Summary|Description|
|------|----|------|-----------|
| POST | `/task` | [Upload files](#task "Upload task file") | Upload the certain deploy files, name: taskID.tar.gz|
| POST | `/command` | [Send commands](#command "Send the testing command") | Tell OCTD which commands need to run|
| GET  | `/result` | [Result](#result "Get the result file") | Get the result file by the path|

###Task
```
POST /task
```

Upload the certain deploy files, name: taskID.tar.gz
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

###command

```
POST /command
```
Tell OCTD which commands need to run. 

**Input**

|Name|Type|Description|
|------|-----|-----------|
| ID | string | The task ID, same with the ID in 'Scheduler' and 'Test Server'|
| Command | string | Command in the deploy.Cmd or run.Cmd in the config.json|

**Response**

```
  { "Status": "ok",
    "Message": "success in running the command"
  }
```

###Result
```
GET /result
```

**Parameters**

| *Name* | *Type* | *Description* |
| -------| ------ | --------- |
| File | string | The result file url. Defined in the config.json in each test case. |
| ID | string | The task ID, same with the ID in 'Scheduler' and 'Test Server' |

**Response**

The whole result file.
```
Server JSON output:
{

        "start":        {
                "connected":    [{
                                "socket":       5,
                                "local_host":   "172.17.0.47",
                                "local_port":   5201,
                                "remote_host":  "172.17.0.49",
                                "remote_port":  59668
                        }],
                "version":      "iperf 3.0.11",
         ......
         ......
         ......
}
```

##Attributes

**Status**
`'OK/Failed'`

**Message**
`The message, especially the error message from the OCTD server.`


