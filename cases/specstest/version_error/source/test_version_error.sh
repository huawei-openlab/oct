#!/bin/bash
cd ./../../source ; runc
if [ $? -eq 0 ]; then
	echo "[NO]        Linuxspec.Spec.Version ==  errValue  failed"
else 
	echo "[Yes]        Linuxspec.Spec.Version == errValue   passed"
fi


