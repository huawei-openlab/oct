# Testing Cases And Framework For Open Container Specifications

The ocp-testing project aims to promote open container project by providing a universal testing framewrk for all the container projects. The testing includes: ocp specification, container function and container performance.

## Framework
### Open Container Pool
    Provide restful API for user who want to query/build/get a container image. 
    The Open Container Pool acts as an agent to deliver requests to different container hubs.
    
### Open Testing Server
    Provide restful API for use who want to use a certain operation system on a certain architect.
    The Open Testing Server acts as an agent to deliver requests to different cluster or IASS platform.

### TestCase Scheduler
    As the main scheduler, the Test Case Scheduler will:
    1) Parse the testing request
    2) Apply hardware resources from the Open Testing Server
    3) Regist container images from the Open Container Pool
    4) Deploy the tesing enviornment
    5) Run the test
    6) Collect and publish the testing report
    
## Testcase specification
Refers to [Cases/README](Cases/README.md) 

## Who should join
- Container project developer
- Operation system distributor
- Hardware company
- IASS provider
- Any container user

## Contributing

The ocp-testing repository is *Apache 2.0* license found in 
the `LICENSE` file of this repository and accepts contributions via GitHub pull requests. 

### Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. 

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
660 York Street, Suite 102,
San Francisco, CA 94110 USA

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

### Email and Chat

### Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for build and test instructions
- Play with the project, submit bugs, submit patches!

### Contribution Flow

### Format of the Commit Message

You just add a line to every git commit message, like this:

    Signed-off-by: Meaglith Ma <maquanyi@huawei.com>

using your real name (sorry, no pseudonyms or anonymous contributions.)

You can add the sign off when creating the git commit via `git commit -s`.
