#!/bin/bash

#start ContainerA
echo -n "Running ContainerA......"
printf "\n"
docker run --name ContainerA -d ubuntu:14.04 /bin/sh
ipaddr=`docker inspect --format '{{ .NetworkSettings.IPAddress}}' ContainerA`
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

docker run --name ContainerB -d ubuntu:14.04 /bin/sh -c "ping $ipaddr > ping-result"
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
