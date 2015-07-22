#!/bin/bash
num_threads=$(cat /proc/cpuinfo | grep processor | wc -l)


cpu_test(){
	for prime_num in 50000 500000  
	do
		sysbench --num-threads=1 --cpu-max-prime=${prime_num} --test=cpu run
		sysbench --num-threads=${num_threads} --cpu-max-prime=${prime_num} --test=cpu run
	done
}


memory_test(){
		sysbench --num-threads=${num_threads} --memory-block-size=1024 --memory-total-size=10G  --test=memory run
		sysbench --num-threads=${num_threads} --memory-block-size=1024 --memory-total-size=100G  --test=memory run
}

fileio_test(){
	for io_option in seqwr seqrd rndwr rndrd 
	do
		sysbench --file-total-size=4096M --test=fileio prepare
		sysbench --file-total-size=4096M  --file-test-mode=${io_option} --test=fileio run
		sysbench --file-total-size=4096M --test=fileio cleanup
		sysbench --file-total-size=8192M --test=fileio prepare
		sysbench --file-total-size=8192M  --file-test-mode=${io_option} --test=fileio run
		sysbench --file-total-size=8192M --test=fileio cleanup
	done
}

cpu_test
memory_test
fileio_test
