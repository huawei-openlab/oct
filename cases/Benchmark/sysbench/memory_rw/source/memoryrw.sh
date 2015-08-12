#!/bin/bash
num_threads=$(cat /proc/cpuinfo | grep processor | wc -l)

memory_test(){
		sysbench --num-threads=${num_threads} --memory-block-size=1024 --memory-total-size=10G  --test=memory run
		sysbench --num-threads=${num_threads} --memory-block-size=1024 --memory-total-size=100G  --test=memory run
}

memory_test
