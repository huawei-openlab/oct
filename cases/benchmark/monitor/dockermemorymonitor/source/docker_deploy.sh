#!/bin/bash

#Fetch build docker image
echo -n "Build docker image"
printf "\n"
docker build -t memoryusage .
printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"

#Start docker container
echo -n "Start docker container"
printf "\n"
docker run -d memoryusage
printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"