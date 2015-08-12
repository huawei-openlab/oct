#!/bin/bash

#Build iperf image
echo -n "Phase1: Building iperf3 image......"
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
printf "\n"

#Running iperf server
echo -n "Phase2: Running iperf server......"
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
printf "\n"

#Beginning udp test
echo -n "Phase3: Start iperf udp testing......"
printf "\n"

docker run -i -t iperf3 -c $ipaddr -u --get-server-output > iperf-udp-result
printf "\n"

b=''
for ((i=0;i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1

    b=#$b
done
printf "\n"
printf "\n"

#Beginning tcp test
echo -n "Phase4: Start iperf tcp testing......"
printf "\n"

docker run -i -t iperf3 -c $ipaddr --get-server-output > iperf-tcp-result
printf "\n"

b=''
for ((i=0;i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1

    b=#$b
done
printf "\n"
printf "\n"

echo -n "Phase5: Check the test result......"
printf "\n"
more iperf-udp-result
more iperf-tcp-result
