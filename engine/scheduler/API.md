#Scheduler
The 'Scheduler' is one of the main components of the whole OCT framework.
It keeps running as a daemon, communicates with the 'Test Case Server'， 'Container service Server' and the 'Test Server'.


- The [configuration](#configs "configuration") file is used for the OCT user to set his/her own configuration.
- The [openAPIs](#apis "APIs") are used for the 'Scheduler' user to monitor the running status.
- The [attributes](#attributes "attributes") are listed at the end of this document.

##Configs
|Key|Type|Description|Example|
|------|----|------| ----- |
| TSurl | string | The url of the Test Server, with the port| "http://localhost:8001" |
| CPurl | string | The url of the Container Service Server, with the port| "http://localhost:8002" |
| RefreshCase | int | Refresh test cases from the Test Case Server every %refreshcase minutes| 30, default to 30|
| Debug | bool | Print the debug information on the screen| true, default to false |

```
{
	"TSurl": "http://localhost:8001",
	"CSSurl": "http://localhost:8002",
	"RefreshCase": 30,
	"Debug": true
}
```


##APIs

|Method|Path|Summary|Description|
|------|----|------|-----------|
| GET | `/task` | [List](#list "list") | List the idle/running/finished tasks. |
| GET | `/task/:ID` | [Details](#details "details") | Fetch the detailed task information. |

###list
```
GET /task
```
List the idle/running/finished tasks.

**Parameters**

| *Name* | *Type* | *Description* |
| -------| ------ | --------- |
| Status |	string | "idle/running/finished". Default: all |
| Page | int | The page number of the tasks (when there are lots of tasks), sort by time. Default: 0 |
| Pagesize | int | The pagesize of the tasks. Default: 10, no more than 100 |

**Response**

```
[
  {
    "ID": "1000",
    "CaseName": "performance/network-latency",
    "Status": "finished",
    "Tested-at": "2015-08-14T16:00:49Z"
  },
  {
    "ID": "1002",
    "CaseName": "function/fake-support",
    "Status": "idle",
    "Tested-at": ""
  }
]
```

###details

```
GET /task/1000
```
Fetch the detailed task information. 

**Response**

```
  {
    "ID": "1000",
    "CaseName": "performance/network-latency",
    "Status": "finished",
    "Report-url": "https://github.com/huawei-openlab/oct/report/****",
    "Tested-at": "2015-08-14T16:00:49Z"
  }
```

##Attributes

**ID**
`The string ID of the testing task.`

**CaseName**
`The name of the test case, same with the API URL for fetching the event.`

**Status**
`'idle/running/finished'`

**Tested_at**
`The timestamp indicating when the testing task occurred (start time).`

**Reported_url**
`The url of the final report.`
