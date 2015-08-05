#!/bin/bash
cd ./../../source/ ; runc
if [ $? -eq 0 ]; then
	echo "[YES]        Linuxspec.Spec.Version == "pre_draft"   passed"
else 
	echo "[NO]        Linuxspec.Spec.Version == "pre_draft"   failed"
fi


