If you just want to try the framework.
The simplest way is to deploy all the micro-services on a simple system and run the testcase.

(The default configure file is already set all the IPs to the local IP address.)

Following these steps:

```
git clone https://github.com/huawei-openlab/oct.git
cd oct/engine
make
cd testserver
./testserver &
cd ../ocitd
./ocitd &
cd ../scheduler/democase
tar czvf ../democase.tar.gz *
cd ..
# choose a testcase.
./scheduler democase.tar.gz

```
The raw result will be store at /tmp/testserver_cache/%(taskID)

For other test cases, please refer to [Test Case Server](tcserver/README.md) to generate them.
