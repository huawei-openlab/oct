# Test Cases and Framework for Open Container Specifications

The oct(open container testing) project aims to promote the [Open Container Initiative](http://www.opencontainers.org/) by providing a universal testing framework for all the container projects. The testing includes: [oci specs](https://github.com/opencontainers/specs), container function and container performance.

## The framework architecture
![Framework](docs/static_files/test_framework.png "Framework")
  * `Open Container Pool` :  
    The Open Container Pool provides RESTful API for user who wants to query/build/get a container image. 
    The Open Container Pool acts as an agent to deliver requests to different container hubs.
    
  * `Open Test Server` :  
    The Open Test Server provides RESTful API for user who wants to use a certain operating system on a certain architecture. 
    The Open Test Server acts as an agent to deliver requests to different cluster or IASS platform.
    
  * `TestCase Scheduler` :  
    As the main scheduler, the Test Case Scheduler will:
    1. Parse the testing request
    2. Apply hardware resources from the Open Test Server
    3. Register container images from the Open Container Pool
    4. Deploy the testing environment
    5. Run the test
    6. Collect and publish the testing report
   
  * `Test Case Server` :  
    The Test Case Server provides RESTful API for user to list/get the test cases.
    It uses the github as the static test case database.
    
## TestCase specification
Refers to [Cases/README](Cases/README.md) 

## Who should join
- Container project developer
- Operating system distributor
- Hardware company
- IASS provider
- PASS provider
- Any container user

### How to involve
If any issues are encountered while using the oct project, several avenues are available for support:
<table>
<tr>
	<th align="left">
	Issue Tracker
	</th>
	<td>
	https://github.com/huawei-openlab/oct/issues
	</td>
</tr>
<tr>
	<th align="left">
	Google Groups
	</th>
	<td>
	https://groups.google.com/forum/#!forum/oci-testing
	</td>
</tr>
</table>

### Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for build and test instructions
- Read the [APIs](engine/API.md) to test each service
- Play with the project, submit bugs, submit patches!
