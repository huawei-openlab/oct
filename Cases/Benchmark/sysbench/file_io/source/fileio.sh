#!/bin/bash

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


fileio_test
