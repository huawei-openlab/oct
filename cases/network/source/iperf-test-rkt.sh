#!/bin/bash

#Build iperf image
echo -n "Phase1: Building iperf3 image......"
printf "\n"
cd iperf-server
docker build -t iperf3 .
cd ..
cd iperf-client-udp
docker build -t iperf3-client-udp .
cd ..
docker build -t iperf3-client-tcp .
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

#Convert docker image to aci
echo -n "Phase2: Transforming docker image iperf3 to iperf3.aci......"
printf "\n"
docker save -o iperf3.docker iperf3
docker save -o iperf3-client-udp.docker iperf3-client-udp
docker save -o iperf3-client-tcp.docker iperf3-client-tcp
docker2aci iperf3.docker
docker2aci iperf3-client-udp.docker
docker2aci iperf3-client-tcp.docker
printf "\n"

#Running iperf server
echo -n "Phase3: Running iperf server......"
printf "\n"

rkt run iperf3.aci
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
echo -n "Phase4: Start iperf udp testing......"
printf "\n"

rkt run iperf3-client-udp.aci
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
echo -n "Phase5: Start iperf tcp testing......"
printf "\n"

rkt run iperf3-client-tcp.aci
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
