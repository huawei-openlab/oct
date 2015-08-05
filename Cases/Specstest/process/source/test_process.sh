#!/bin/bash
cd ./../../source ; runc
if [ $? -eq 0 ]; then
	echo "[YES]        Linuxspec.Spec.Process    passed"
else 
	echo "[NO]        Linuxspec.Spec.Process    failed"
fi


