#Test Case Server
The 'Test Case Server' provides all the test cases which stored at github.

- The [configuration](#configs "configuration") file is used for the OCT users to set his/her own configuration.
- The [openAPIs](#apis "APIs") are used for the developers to integrate the 'Test Case Server' with the 'Scheduler'.
- The [attributes](#attributes "attributes") are listed at the end of this document.

##Configs
|Key|Type|Description|Example|
|------|----|------| ----- |
| GitRepo | string | The git repo url.| "https://github.com/huawei-openlab/oct.git" |
| CaseDir | string | The case dir of the git repo.| "cases" |
| Group | []string | The group names, also means the sub directory| ["Performance"]|
| CacheDir | string | The cache dir where the repo stored.| "/tmp/tcserver_cache" |
| Debug | bool | Print the debug information on the screen| true, default to false |

```
{
	"GitRepo":  "https://github.com/huawei-openlab/oct.git",
        "CaseDir":  "Cases",
	"Group": ["Network", "Benchmark", "Specstest"],
	"CacheDir": "/tmp/tcserver_cache",
	"Port":  8011
}
```


##APIs

|Method|Path|Summary|Description|
|------|----|------|-----------|
| GET | `/case` | [List](#list "list") | List the cases. |
| GET | `/case/:ID` | [Details](#details "details") | Fetch the case files. %caseid.tar.gz |
| GET | `/case/:ID/report` | [Report] (#report "report") | Get the case report|
| POST| `/case` | [Refresh](#refresh "refresh") | Refresh the cases. |

###list
```
GET /case
```
List the idle/tested test cases.

**Parameters**

| *Name* | *Type* | *Description* |
| -------| ------ | --------- |
| Status |	string | "idle/tested". Default: all |
| Page | int | The page number of the test cases, sort by time. Default: 0 |
| Pagesize | int | The pagesize of the test cases. Default: 10, no more than 100 |

**Response**

```
[
  {
    "ID": "10100",
    "CaseName": "performance/network-latency",
    "Status": "tested"
  },
  {
    "ID": "10102",
    "CaseName": "function/fake-support",
    "Status": "idle"
  }
]
```

###details

```
GET /case/10100
```
Fetch the case files. 

**Response**

The whole %caseid.tar.gz file.

###report

```
GET /case/10100/report
```
Fetch the case report. 

**Response**

The whole %caseid-report file.


###refresh
```
POST /case
```
Refresh the test cases.

**Response**

```
 {
    "Status": "OK",
    "Message": ""
 }
```

##Attributes

**ID**
`The string ID of the test case.`

**CaseName**
`The name of the test case, same with the API URL for fetching the event.`

**Status**
In `List` : `'idle/tested'`

In `Refresh` : `'OK/Failed'`



