##The testcase should be written in JSON format, including five parts:
```
  "Description" part, "Requires", "Deploy", "Run" and "Collect". 
```

- We provide a 'casevalidator' for case writers, Read the [HowTO](../engine/tools/casevalidator/HowTO.md) to check if your case was valid.
- The cases in this oct repo is good to use, if you want to try any of them, please read the [Test Case Server](../engine/tcserver/README.md).


### "Description" part
The case developer should fill in the following informations:
  "Name", "Owner", "Version", "Summary", "Licence", "Group", "URL", "Description", "Source". 
(Yes, looks like rpm spec)
```
{"Name": "Test001", 
 "Owner", "developer001@test.com", 
 "Summary": "demo test",
 "License": "Apache 2.0" //default
 "Group": "Performance/Network",  // The first category group is Performance/Spec/Function
                                  // The 'Group' is used to category each test case, 
                                  // we assume/hope that there will be 'lots of' test cases, 
                                  // so search one by the group
 "URL": "github://test",
 "Description": "This is the detail description"
}
```

### "Requires"
The case developer should define the type, the hardware requirements and the ammount of the host opertation systems.
Also, he/she should define the container image (by dockerfile for example), so the backend server will find/build an image.
```
 "Requires": [
            {
                "Class": "operationOS",  //class, not an object
                "Type": "os",  // only 'os' and 'container'
                "Distribution": "Ubuntu",  //case insensitive
                "Version": "14.02"
                "Resource": {
                        "CPU": 1,
                        "Memory": "2GB",
                        "Disk": "100G"
                }
                "amount": 2  //how many hostOS do we need, no need if the 'type' is container
            }
            {
                "Class": "ImageA": 
                "Type": "container",
                "Version": "dockerV1.0"
                "Files": ["a.dockerfile"]
            }
        ],
```

### "Deploys"
The case developer should tell the host operation systems and the containers 
which commands should be used to deploy the test.
The commands could be used directly or wrapped by a script. 
By default, after runing all the command, the test system will continue to the 'Run' part.
If the developer want to install extra package, he/she can put the related commands here. 
```
  "Deploys": [
            {
                "Object": "hostA",
                "Class": "OperationOS",
                "Files": ["./source/b.sh", "./source/a.dockerfile"],
                 "Containers": [
                        {
                            "Object": "DockerA",
                            "Class": "ImageA",
                            "Cmd": ""
                        }
                  ]
                
            },
            {
                "Object": "hostB",
                "Class":  "OperationOS",  //difference object, same class
                "Containers": [
                        {
                            "Object": "DockerB",
                            "Class": "ImageA",
                            "Cmd": "systemctl start sshd"
                            }
                        }
                    ]
            }
        ],
```         

### "Run"
The case developer should tell the host operation systems and the containers
which commands should be used to run the test. The commands could be used 
directly or wrapped by a script. By default, after ruuning all the command, the test 
system will continue to the 'Collect' part.
```
"Run": [
            { "Object": "HostA",
              "Cmd": ""
            }, 
            { "Object": "DockerA",
              "Cmd" : "ping -c 1 -s 1500 192.168.10.10"
            }
      ]
```

### "Collects"
The case developer should tell the host operation systems and/or the containers 
if there was any output file.
```
"Collects": [
        {"Object": "HostA",
          "Files": ["./source/output.json"]  //this file will be return to our framework as the output
        }
  ]
```
