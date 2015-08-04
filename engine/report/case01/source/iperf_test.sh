#!/bin/bash

#Build iperf image
echo -n "Building iperf3 image"
printf "\n"
docker build -t iperf3 .
printf "\n"
b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"

#Running iperf server
echo -n "Running iperf server"
printf "\n"
docker run --name iperf_Server -d iperf3 -s -J
ipaddr=`docker inspect --format '{{ .NetworkSettings.IPAddress}}' iperf_Server`

printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"

#Beginning udp test
echo -n "Start iperf udp test"
printf "\n"
#docker run  -i -t iperf3 -c $ipaddr -u --get-server-output > iperf-udp-result.json 
docker run  -i -t iperf3 -c $ipaddr -u --get-server-output 1 > iperf-udp-result.json 
printf "\n"
b=''
for ((i=0;i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1

    b=#$b
done
echo -n "Finished iperf udp test "
printf "\n"

#Beginning tcp test
echo -n "Start iperf tcp test"
printf "\n"
#docker run -i -t iperf3 -c $ipaddr --get-server-output > iperf-tcp-result.json
docker run -i -t iperf3 -c $ipaddr --get-server-output 1 > iperf-tcp-result.json
printf "\n"
b=''
for ((i=0;i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1

    b=#$b
done
echo -n "Finished iperf tcp test"
printf "\n"

ocitd -s "Finish"
