cd /oct/cases/benchmark/monitor
mkdir tmp_test
cd tmp_test
cp -fr ../dockercpumonitor .
cp -fr ../source .
tar czvf ../dockercpumonitor.tar.gz *
cd ..
rm -fr tmp_test
