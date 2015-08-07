#!/bin/bash

#Download ubuntu docker image
echo -n "Pulling ubuntu 14.04 image......"
printf "\n"
docker pull ubuntu:14:04
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
echo -n "Transforming docker image ubuntu to ubuntu.aci......"
printf "\n"
docker save -o ubuntu.docker ubuntu:14.04
docker2aci ubuntu.docker
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

#start ContainerA
echo -n "Running ContainerA......"
printf "\n"
rkt run ubuntu.aci ----insecure-skip-verify
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

#Start Container B and test
echo -n "Running ContainerB and ping Container A......"
printf "\n"

rkt run ubuntu.aci ----insecure-skip-verify
ping $ipaddr
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
