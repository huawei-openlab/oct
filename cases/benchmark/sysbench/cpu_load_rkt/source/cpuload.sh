#!/bin/bash
num_threads=$(cat /proc/cpuinfo | grep processor | wc -l)


cpu_test(){
	for prime_num in 50000 500000  
	do
		sysbench --num-threads=1 --cpu-max-prime=${prime_num} --test=cpu run
		sysbench --num-threads=${num_threads} --cpu-max-prime=${prime_num} --test=cpu run
	done
}


cpu_test

