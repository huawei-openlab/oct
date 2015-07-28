## The testcases should be written in JSON format, Each testcase should conclude a souce fold and a JSON file

### JSON file
The JSON file should be called by the test scheduler to do test jobs, include serveral parts:
```
{"Description": {}, "Requires": {}, "Deploy":{}, "Run": {}, "Collect": {}}
```
For more information, refers to [ocp-testing/README.md](./../../README.md).

### source folder
The source folder contains souce code needed by testcases.
For example ,the testcase Benchmark-monitor-cpumonitor-totalusage-docker.json have three souce code file:
Dockerfile, cpu_total_usage.go and deadloop.sh, you can find them in 
```
ocp-testing/Cases/Benchmark/monitor/cpumonitor/source/ 
```
Meanwhile, the Benchmark-monitor-cpumonitor-totalusage-docker.json tells the scheduler how to exec them:
```
"root": {
        "Name": "Benchmark-Test-Monitor-docker",
        "Summary": "test with cadvisor tool",
        "Owner": "linzhinan@huawei.com",
        "Details": "Monitor cpu totalusage via cadvisor",
        "Group": "Performance/Benchmark/monitor/cpumonitor/",
        "License": "Apache 2.0",
        "Source": [
                    "./source/Dockerfile",
	       "./source/cpu_total_usage.go",
                    "./source/deadloop.sh"
        ],
```
It tells the scheduler we have three source code file ,cpu_total_usage.go and deadloop.sh.
```
 "Run": [
            {
                "hostA": {
                    "run": {
                        "type": "cmd",
                        "value": "cd ./source/",
                        "type": "cmd",
                        "value": "go build cpu_total_usage.go",
                        "cmd": "cmd",
                        "type": "sudo ./cpu_total_usage"
                    }
                    
                }
            },
```
It tells the scheduler to compile and run the cpu_total_usage test client.
```
"containers": [
                    {
                        "name": "DockerA",
                        "class": "ImageA",
                        "setup": {
                            "type": "cmd",
                            "value": "cd /",
                            "type": "cmd",
                            "value": "./deadloop.sh"
                        }
                    }
                ]
```
Because of the close to zero cpuload in none task container, the testcase installed a deadloop.sh to the container to ensurence it can show the result of the monitor easily.  

For more information, refers to [ocp-testing/README.md](./../../README.md).

## The Benchmark testcases conclude two parts, monitor benchmark and sysbench benchmark

### monitor benchmark
It use the cadvisor as a server, cadvisor is developped by google, it can montor the containers status and export the data to the websit, 
for more information, refers to https://github.com/google/cadvisor
You can find the called code in JSON file, for example, in Benchmark-monitor-cpumonitor-totalusage-docker.json,
```
      "Deploy": [
            {
                "object": "hostA",
                "class": "OperationOS",
                "setup": {
                    "type": "cmd",
                    "value": "cd /tmp",
                    "type": "cmd",
                    "value": "git clone https://github.com/google/cadvisor.git",
                    "type": "cmd",
                    "value": "godep go build .",
                    "type": "cmd",
                    "value": "echo hostpasswd | sudo -S ./cadvisor"
                },
```
Each go code file in source folder is a cadvisor client,  the testcase use cadvisor client api to collect metrics from cadvisor server.

### sysbench benchmach
It use the sysbench tool to run pressure test and collect the cpu load mtrics.
Sysbench is installed in the container, to the docker container, u can find the installed action in
``` 
ocp-testing/Cases/Benchmark/sysbench/xxx/source/Dockerfile,
```
```
# Install sysbench
RUN sudo apt-get update
RUN sudo apt-get install sysbench -y
# Install sysbench
RUN sudo apt-get update
RUN sudo apt-get install sysbench -y
```