The testcase should be written in JSON format, including five parts: Description,
Requires, Deploy, Run, Collect.
{"Description": {}, "Requires": {}, "Deploy":{}, "Run": {}, "Collect": {}}

Description
We should fill in the case name, the maintainer, the license and the summary.
The group is used to category each test case, in the future, we assume that there 

will be 'lots of' test cases, so search one by the group --- change and run it.


Requires
In this part, the case developer should define the type, the hardware requirement and 

ammount of the host opertation systems. Also, he/she should define the container 

image (by dockerfile for example).

Deploy
In this part, the case developer should tell the host operation systems and 

containers which commands should be used to deploy the test. The commands could be 

used directly or wrapped by a script. By default, after ruuning all the command, the 

test system will start to the 'Run' part.

If the developer want to install extra package, he/she can put the related command 

here. 

Run
In this part, the case developer should tell the host operation systems and the 

containers which commands should be used to run the test. The commands could be used 

directly or wrapped by a script. By default, after ruuning all the command, the test 

system will start to the 'Collect' part.

on the host os and  extra packages should be installedcommands 


'Performance/Network'
