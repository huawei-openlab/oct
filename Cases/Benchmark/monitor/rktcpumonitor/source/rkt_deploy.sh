#!/bin/bash

#Fetch build docker image
echo -n "Fetch rkt image of coreos.com/etcd:v2.0.4 "
rkt fetch  --insecure-skip-verify coreos.com/etcd:v2.0.4
printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done

#Start contianer etcd:v2.0.4
rkt run --mds-register=false --local --interactive --debug coreos.com/etcd:v2.0.4 &
printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"