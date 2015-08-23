#!/bin/bash
runc
if [ $? -eq 0 ]; then
	echo  "{\"Linuxspec.Mounts.Filesystem\":{\"$1\":\"success\"}}"  > /tmp/testtool/mount_fstype_out.json
	# echo "[NO]   Filesystem $1 Mount Unsupported " >
else 
	echo  "{\"Linuxspec.Mounts.Filesystem\":{\"$1\":\"failed\"}}" > /tmp/testtool/mount_fstype_out.json
	# echo "[Yes]   Filesystem $1 Mount Supported "
fi


