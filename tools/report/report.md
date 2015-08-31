## Network-Test_002
[Test Case](#testcase) Test network bandwidth between different containers

```
Owner: chengtiesheng@huawei.com
License: Apache 2.0
Group: Performance/Network
```

The case has 1 operation system(s):


'hostA' has 2 container(s) deployed.



The defailed information is listed as below:

| *OS Name* | *Distribution* | *Resource* | *Container*| *Deploy/Testing Command* |
| -------| ------ | --------- | -------- | --------|
|hostA|CentOS7|CPU 1, Memory 2GB, Disk 100G|iperf_server(iperf3), iperf_client(iperf3)| ./iperf_test.sh|


The defailed information of each container type is listed as below:

| *Container Type* | *Distribution* | *Container File* |
| -------| ------ | ------- |
|iperf3|Docker1|[Dockerfile](#dockerfile)|


After running the `Command` in each OS and container, we get 2 result(s).


* [iperf-tcp-result.json](#iperf-tcp-result)

* [iperf-udp-result.json](#iperf-udp-result)


###testCase

```
{
        "Name": "Network-Test_002",
        "Summary": "network bandwidth",
        "Version": "0.1",
        "Owner": "chengtiesheng@huawei.com",
        "Description": "Test network bandwidth between different containers",
        "Group": "Performance/Network",
        "License": "Apache 2.0",
        "Explains": "The Sources here is useless",
        "Sources": [
            "./source/Dockerfile",
            "./source/iperf_test.sh"
        ],
        "Requires": [
            {
                    "Class": "operOS",
                    "Type": "os",
            		    "Distribution": "CentOS",
		                "Version": "7",
                    "Resource": {
                        "CPU": 1,
                        "Memory": "2GB",
                        "Disk": "100G"
                    }
            },
            {
                    "Class": "iperf3",
		    "Type": "container",
		    "Distribution": "Docker",
                    "Version": "1",
		    "Files": ["./source/Dockerfile"]
            }
        ],
        "Deploys": [
            {
                "Object": "hostA",
                "class": "operOS",
		"Files": ["./source/iperf_test.sh", "./source/Dockerfile"],
                "Containers": [
                        {
                            "Object": "iperf_server",
                            "Class": "iperf3"
                        },
                        {
                            "Object": "iperf_client",
                            "Class": "iperf3"
                        }
                    ],
		"Cmd": "./iperf_test.sh"
            }
	],
        "Collects": [
            {
                "Object": "hostA",
		"Files": ["./source/iperf-tcp-result.json", "./source/iperf-udp-result.json"]
            }
        ]
}

```



###dockerfile

```
FROM ubuntu:14.04

MAINTAINER Tiesheng Cheng <chengtiesheng@huawei.com>

ENV VERSION=3.0.11 \
    BUILD_DEPS="gcc libc6-dev make curl"

RUN \ 
apt-get update && \
    apt-get -y install --no-install-recommends $BUILD_DEPS && \
    curl -sL http://downloads.es.net/pub/iperf/iperf-${VERSION}.tar.gz | tar xzf - -C /tmp --strip 1 && \
    (cd /tmp; ./configure) && \
    make -C /tmp && \
    make -C /tmp install && \
    apt-get --purge autoremove -y $BUILD_DEPS && \
    apt-get clean && \
    rm -rf /var/cache/* /var/lib/apt/lists/* /tmp/*

ENTRYPOINT ["/usr/local/bin/iperf3"]

```




###iperf-tcp-result

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
                "system_info":  "Linux 7144f18e1e5d 3.10.0-229.7.2.el7.x86_64 #1 SMP Tue Jun 23 22:06:11 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux\n",
                "timestamp":    {
                        "time": "Mon, 03 Aug 2015 17:30:55 GMT",
                        "timesecs":     1438623055
                },
                "accepted_connection":  {
                        "host": "172.17.0.49",
                        "port": 59667
                },
                "cookie":       "361e905f5b10.1438623055.947216.067c4",
                "tcp_mss_default":      1448,
                "test_start":   {
                        "protocol":     "TCP",
                        "num_streams":  1,
                        "blksize":      131072,
                        "omit": 0,
                        "duration":     10,
                        "bytes":        0,
                        "blocks":       0,
                        "reverse":      0
                }
        },
        "intervals":    [{
                        "streams":      [{
                                        "socket":       5,
                                        "start":        0,
                                        "end":  1.00003,
                                        "seconds":      1.00003,
                                        "bytes":        2312923104,
                                        "bits_per_second":      18502900000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        0,
                                "end":  1.00003,
                                "seconds":      1.00003,
                                "bytes":        2312923104,
                                "bits_per_second":      18502900000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        1.00003,
                                        "end":  2.00001,
                                        "seconds":      0.999987,
                                        "bytes":        2408972288,
                                        "bits_per_second":      19272000000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        1.00003,
                                "end":  2.00001,
                                "seconds":      0.999987,
                                "bytes":        2408972288,
                                "bits_per_second":      19272000000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        2.00001,
                                        "end":  3.00001,
                                        "seconds":      1,
                                        "bytes":        2411724800,
                                        "bits_per_second":      19293800000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        2.00001,
                                "end":  3.00001,
                                "seconds":      1,
                                "bytes":        2411724800,
                                "bits_per_second":      19293800000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        3.00001,
                                        "end":  4.00002,
                                        "seconds":      1.00001,
                                        "bytes":        2405761400,
                                        "bits_per_second":      19245900000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        3.00001,
                                "end":  4.00002,
                                "seconds":      1.00001,
                                "bytes":        2405761400,
                                "bits_per_second":      19245900000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        4.00002,
                                        "end":  5.00001,
                                        "seconds":      0.999988,
                                        "bytes":        2384658056,
                                        "bits_per_second":      19077500000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        4.00002,
                                "end":  5.00001,
                                "seconds":      0.999988,
                                "bytes":        2384658056,
                                "bits_per_second":      19077500000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        5.00001,
                                        "end":  6.00001,
                                        "seconds":      1,
                                        "bytes":        2386821120,
                                        "bits_per_second":      19094500000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        5.00001,
                                "end":  6.00001,
                                "seconds":      1,
                                "bytes":        2386821120,
                                "bits_per_second":      19094500000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        6.00001,
                                        "end":  7.00003,
                                        "seconds":      1.00001,
                                        "bytes":        2388131840,
                                        "bits_per_second":      19104800000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        6.00001,
                                "end":  7.00003,
                                "seconds":      1.00001,
                                "bytes":        2388131840,
                                "bits_per_second":      19104800000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        7.00003,
                                        "end":  8.00001,
                                        "seconds":      0.999978,
                                        "bytes":        2388656128,
                                        "bits_per_second":      19109700000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        7.00003,
                                "end":  8.00001,
                                "seconds":      0.999978,
                                "bytes":        2388656128,
                                "bits_per_second":      19109700000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        8.00001,
                                        "end":  9.00002,
                                        "seconds":      1.00002,
                                        "bytes":        2391015424,
                                        "bits_per_second":      19127700000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        8.00001,
                                "end":  9.00002,
                                "seconds":      1.00002,
                                "bytes":        2391015424,
                                "bits_per_second":      19127700000,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        9.00002,
                                        "end":  10,
                                        "seconds":      0.999985,
                                        "bytes":        2389835776,
                                        "bits_per_second":      19119000000,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        9.00002,
                                "end":  10,
                                "seconds":      0.999985,
                                "bytes":        2389835776,
                                "bits_per_second":      19119000000,
                                "omitted":      false
                        }
                }],
        "end":  {
        }
}

```


###iperf-udp-result

```
Server JSON output:
{
        "start":        {
                "connected":    [{
                                "socket":       5,
                                "local_host":   "172.17.0.47",
                                "local_port":   5201,
                                "remote_host":  "172.17.0.48",
                                "remote_port":  59576
                        }],
                "version":      "iperf 3.0.11",
                "system_info":  "Linux 7144f18e1e5d 3.10.0-229.7.2.el7.x86_64 #1 SMP Tue Jun 23 22:06:11 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux\n",
                "timestamp":    {
                        "time": "Mon, 03 Aug 2015 17:30:40 GMT",
                        "timesecs":     1438623040
                },
                "accepted_connection":  {
                        "host": "172.17.0.48",
                        "port": 33435
                },
                "cookie":       "d9c0e27fc62a.1438623040.417897.57020",
                "test_start":   {
                        "protocol":     "UDP",
                        "num_streams":  1,
                        "blksize":      8192,
                        "omit": 0,
                        "duration":     10,
                        "bytes":        0,
                        "blocks":       0,
                        "reverse":      0
                }
        },
        "intervals":    [{
                        "streams":      [{
                                        "socket":       5,
                                        "start":        0,
                                        "end":  1.00008,
                                        "seconds":      1.00008,
                                        "bytes":        122880,
                                        "bits_per_second":      982961,
                                        "jitter_ms":    0.00995991,
                                        "lost_packets": 0,
                                        "packets":      15,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        0,
                                "end":  1.00008,
                                "seconds":      1.00008,
                                "bytes":        122880,
                                "bits_per_second":      982961,
                                "jitter_ms":    0.00995991,
                                "lost_packets": 0,
                                "packets":      15,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        1.00008,
                                        "end":  2.00009,
                                        "seconds":      1.00001,
                                        "bytes":        131072,
                                        "bits_per_second":      1048560,
                                        "jitter_ms":    0.0145827,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        1.00008,
                                "end":  2.00009,
                                "seconds":      1.00001,
                                "bytes":        131072,
                                "bits_per_second":      1048560,
                                "jitter_ms":    0.0145827,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        2.00009,
                                        "end":  3.00011,
                                        "seconds":      1.00001,
                                        "bytes":        131072,
                                        "bits_per_second":      1048560,
                                        "jitter_ms":    0.0153738,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        2.00009,
                                "end":  3.00011,
                                "seconds":      1.00001,
                                "bytes":        131072,
                                "bits_per_second":      1048560,
                                "jitter_ms":    0.0153738,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        3.00011,
                                        "end":  4.00009,
                                        "seconds":      0.999988,
                                        "bytes":        131072,
                                        "bits_per_second":      1048590,
                                        "jitter_ms":    0.0174913,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        3.00011,
                                "end":  4.00009,
                                "seconds":      0.999988,
                                "bytes":        131072,
                                "bits_per_second":      1048590,
                                "jitter_ms":    0.0174913,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        4.00009,
                                        "end":  5.00009,
                                        "seconds":      1,
                                        "bytes":        131072,
                                        "bits_per_second":      1048580,
                                        "jitter_ms":    0.0166612,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        4.00009,
                                "end":  5.00009,
                                "seconds":      1,
                                "bytes":        131072,
                                "bits_per_second":      1048580,
                                "jitter_ms":    0.0166612,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        5.00009,
                                        "end":  6.00009,
                                        "seconds":      0.999998,
                                        "bytes":        131072,
                                        "bits_per_second":      1048580,
                                        "jitter_ms":    0.0158326,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        5.00009,
                                "end":  6.00009,
                                "seconds":      0.999998,
                                "bytes":        131072,
                                "bits_per_second":      1048580,
                                "jitter_ms":    0.0158326,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        6.00009,
                                        "end":  7.00011,
                                        "seconds":      1.00002,
                                        "bytes":        131072,
                                        "bits_per_second":      1048560,
                                        "jitter_ms":    0.0167468,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        6.00009,
                                "end":  7.00011,
                                "seconds":      1.00002,
                                "bytes":        131072,
                                "bits_per_second":      1048560,
                                "jitter_ms":    0.0167468,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        7.00011,
                                        "end":  8.00009,
                                        "seconds":      0.999987,
                                        "bytes":        131072,
                                        "bits_per_second":      1048590,
                                        "jitter_ms":    0.0157667,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        7.00011,
                                "end":  8.00009,
                                "seconds":      0.999987,
                                "bytes":        131072,
                                "bits_per_second":      1048590,
                                "jitter_ms":    0.0157667,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                        "socket":       5,
                                        "start":        8.00009,
                                        "end":  9.00009,
                                        "seconds":      0.999999,
                                        "bytes":        131072,
                                        "bits_per_second":      1048580,
                                        "jitter_ms":    0.0164364,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        8.00009,
                                "end":  9.00009,
                                "seconds":      0.999999,
                                "bytes":        131072,
                                "bits_per_second":      1048580,
                                "jitter_ms":    0.0164364,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }, {
                        "streams":      [{
                                       "socket":       5,
                                        "start":        9.00009,
                                        "end":  10.0001,
                                        "seconds":      1,
                                        "bytes":        131072,
                                        "bits_per_second":      1048576,
                                        "jitter_ms":    0.0170041,
                                        "lost_packets": 0,
                                        "packets":      16,
                                        "lost_percent": 0,
                                        "omitted":      false
                                }],
                        "sum":  {
                                "start":        9.00009,
                                "end":  10.0001,
                                "seconds":      1,
                                "bytes":        131072,
                                "bits_per_second":      1048576,
                                "jitter_ms":    0.0170041,
                                "lost_packets": 0,
                                "packets":      16,
                                "lost_percent": 0,
                                "omitted":      false
                        }
                }],
        "end":  {
        }
}

```

