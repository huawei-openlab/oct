If you just want to try the framework.
The simplest way is to deploy all the micro-services on a simple system and run the testcase.

(The default configure file is already set all the IPs to the local IP address.)

Following these steps:
- Clone this project
- Go to scheduler/testserver/iocitd/containerpool directory
  build the related .go file.
- Run testserver/iocitd/containerpool as daemons
- Go to scheduler directory, choose a testcase.

  ``` Example
   ./scheduler case01/Network-iperf.json
  ```
  
  The raw result will be store at /tmp/test_scheduler_result
