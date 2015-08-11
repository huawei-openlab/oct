#!/bin/bash
cd ./../../source ; runc 
if [ $? -eq 0 ]; then
	echo "[NO]   Filesystem $1 Mount Unsupported "
else 
	echo "[Yes]   Filesystem $1 Mount Supported "
fi


