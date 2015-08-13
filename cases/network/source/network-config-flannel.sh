#!/bin/sh
# This code should (try to) follow Google's Shell Style Guide
# (https://google-styleguide.googlecode.com/svn/trunk/shell.xml)
curl -L https://github.com/coreos/etcd/releases/download/v2.0.9/etcd-v2.0.9-linux-amd64.tar.gz -o etcd-v2.0.9-linux-amd64.tar.gz
git clone
apt-get install linux-libc-dev
cd flannel; ./build

tar -zxvf etcd-v2.0.9-linux-amd64.tar.gz
cd etcd-v2.0.9-linux-amd64/
./etcd --listen-client-urls=http://0.0.0.0:4001 --listen-peer-urls=http://0.0.0.0:7001 &
./etcdctl mk /coreos.com/network/config '{"Network":"192.168.0.0/16"}'

ip link set dev docker0 down
brctl delbr docker0
cd flannel/bin/
./flanneld -iface=eth0 & 
source /run/flannel/subnet.env
docker -d --bip=${FLANNEL_SUBNET} --mtu=${FLANNEL_MTU} &
