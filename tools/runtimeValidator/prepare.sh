#!/bin/bash
set -e

TESTROOT=./bundles/$1
rm -rf ${TESTROOT}
mkdir -p ${TESTROOT}
tar -xf  rootfs.tar.gz -C ${TESTROOT}

cp ./plugins/runtimetest ${TESTROOT}

