#!/bin/bash
echo "Going to prepare testcase:"
echo $2
rm -rf tmp_dir
mkdir tmp_dir
cd tmp_dir
cp -rf ./../$2 .
cp -rf ./../source .
tar czvf ../$2.tar.gz *

cd ..
rm -rf tmp_dir
cp $2.tar.gz ./../../engine/scheduler/

cd ./../../engine/scheduler
./scheduler ./$2.tar.gz
