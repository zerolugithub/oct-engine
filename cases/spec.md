#OCT Case Specifications
The testcase for OCT should comply with the following specifications:
- [Case Bundle](#case-bundle)
  - [Directory Layout](#directory-layout)
- [Case Configuration](#case-configuration)
  - [Description](#case-description)
  - [Requires](#case-requires)
  - [Deploys](#case-deploys)
  - [Run](#case-run)
  - [Collects](#case-result)

##Case Bundle
The term `bundle` is inspired by [OCI-bundle](https://github.com/opencontainers/specs/blob/master/bundle.md).
Case Bundle is a directory organized in a certain way, containing one configuration file and all the necessary testing scripts. The spec does not specify how to name, tar or transfer a case bundle, these will be considered as the '[implementation](#implementation)'. 

###Directory Layout
The layout of a standard case bundle is like this, all the testing scripts with/without sub-directoies should be put in the 'source' directory.
```
config.json
source/
```

##Case Configuration
The testcase should be written in JSON format, including five parts:
```
  "Description" part, "Requires", "Deploys", "Run" and "Collects". 
```

###Case Description
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

###Case Requires
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

###Case Deploys
The case developer should tell the host operation systems and the containers 
which commands should be used to deploy the test.
The commands could be used directly or wrapped by a script. 
By default, after runing all the command, the test system will continue to the 'Run' part.
If the developer want to install extra package, he/she can put the related commands here. 
The default `CWD` in `Cmd` is the 'source' directory.
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

###Case Run
The case developer should tell the host operation systems and the containers
which commands should be used to run the test. The commands could be used 
directly or wrapped by a script. By default, after running all the command, the test 
system will continue to the 'Collects' part.
The default `CWD` in `Cmd` is the 'source' directory.
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

###Case Result
The case developer should tell the host operation systems and/or the containers 
if there was any output file.
The default `CWD` in `Files` is the 'source' directory.
```
"Collects": [
        {"Object": "HostA",
          "Files": ["output.json"]  //this file will be return to our framework as the output
        }
  ]
```
