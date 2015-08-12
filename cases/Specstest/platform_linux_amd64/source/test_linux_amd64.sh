#!/bin/bash
cd ./../../source/ ; runc
if [ $? -eq 0 ]; then
	echo "[YES]        linuxspec.Spec.Platform.OS = "linux"     linuxspec.Spec.Platform.Arch = "amd64"   passed"
else 
	echo "[NO]        linuxspec.Spec.Platform.OS = "linux"     linuxspec.Spec.Platform.Arch = "amd64"   failed"
fi


