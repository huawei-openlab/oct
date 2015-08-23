#!/bin/bash
cd ./../../engine
echo "Kill all of the deamon task of oct ..."
echo "\n"
kill -9 `ps -ef|grep testserver|grep -v grep|awk '{print $2}'`
kill -9 `ps -ef|grep ocitd|grep -v grep|awk '{print $2}'`
kill -9 `ps -ef|grep runc|grep -v grep|awk '{print $2}'`

echo "Rm all of the oci runc resource ..."
echo "\n"
rm -r /run/oci/source/

echo "Do make clean ...\n"

make clean

echo "Start make the oct framework..."
make

echo "Start testserver deamon ..."
echo "\n"

cd testserver
./testserver &

echo "Start ocitd deamon"
echo "\n"

cd ../ocitd
./ocitd &
