#!/bin/sh
# This code should (try to) follow Google's Shell Style Guide
# (https://google-styleguide.googlecode.com/svn/trunk/shell.xml)
git clone https://github.com/jpetazzo/pipework
cd pipework/
cp pipework /usr/bin

pipework br0 $container 172.28.1.2/24@172.28.1.1
brctl addif br0 eth0
ip addr add 172.28.1.1/24 dev br0
