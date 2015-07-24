TODO-List (by seqence)

1. testcase server
1.1 parse & validation
    find the mandatory area
1.2 send request to test-server and container-server
1.3 add config file to tell the test-server daemon and container pool daemon address

2. test-server-daemon
2.1 monitor the request from testcase server
2.2 deliver to the underling HA/openStack...
2.3 return resource id/id-list to testcase server

3. container pool daemon
3.1 monitor the request from testcase server
3.2 deliver to the undeerling container hub
3.3 return resource id/id-list to testcase server

4. testcase server send new testcase order (with both id/id-list) to test-server-daemon
    to deploy the test

# go-hack
start to use golang
