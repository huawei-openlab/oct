#!/bin/bash
cd ./../../source ; runc
if [ $? -eq 0 ]; then
	echo "[NO]        Filesystem Mount Unsupported " 
else 
	echo "[Yes]        Filesystem Mount Supported "
fi


