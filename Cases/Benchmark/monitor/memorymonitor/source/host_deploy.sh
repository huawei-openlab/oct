#!/bin/bash

pre_dir=`pwd`
#Fetch cadvisor source code
echo -n "Fetch cadvisor from github"

printf "\n"
mkdir -p $GOPATH/src/github.com/google/
cd $GOPATH/src/github.com/google/
git clone https://github.com/google/cadvisor.git
printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"

#Build cadvisor
echo -n "Build cadvisor"
printf "\n"
cd $GOPATH/src/github.com/google/cadvisor
godep go build .
printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"

#Run cadvisor in host
echo -n "Start cadvisor"
printf "\n"
./cadvisor &
printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"

#Compile  memorymonitor on host
echo  "Compile memorymonitor.go"
cd $pre_dir
if [ -n $2 ] ;
	then
	go build $2
fi
#go build memorymonitor.go
printf "\n"

b=''
for ((i=0;$i<=100;i+=2))
do
    printf "progress:[%-50s]%d%%\r" $b $i
    sleep 0.1
    b=#$b
done
printf "\n"