# Testing Cases And Framework For Open Container Specifications

The oct(open container testing) project aims to promote the open container project by providing a universal testing framewrk for all the container projects. The testing includes: ocp specification, container function and container performance.

## The framework architecture
![Framework](docs/static_files/test_framework.png "Framework")
  * `Open Container Pool` :  
    The open contaner Pool provides restful API for user who want to query/build/get a container image. 
    The Open Container Pool acts as an agent to deliver requests to different container hubs.
    
  * `Open Test Server` :  
    The Open Test Server provides restful API for use who want to use a certain operation system on a certain architecture. 
    The Open Test Server acts as an agent to deliver requests to different cluster or IASS platform.
    
  * `TestCase Scheduler` :  
    As the main scheduler, the Test Case Scheduler will:
    1. Parse the testing request
    2. Apply hardware resources from the Open Test Server
    3. Regist container images from the Open Container Pool
    4. Deploy the tesing enviornment
    5. Run the test
    6. Collect and publish the testing report
   
  * `Test Case Server` :  
    The Test Case Server provides restful API for user to list/get the test cases.
    It uses the github as the static test case database.
    
## Testcase specification
Refers to [Cases/README](Cases/README.md) 

## Who should join
- Container project developer
- Operation system distributor
- Hardware company
- IASS provider
- Any container user

### How to involve
- Mailing list [Join](https://groups.google.com/forum/#!forum/oci-testing)

### Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for build and test instructions
- Read the [APIs](engine/API.md) to test each service
- Play with the project, submit bugs, submit patches!
