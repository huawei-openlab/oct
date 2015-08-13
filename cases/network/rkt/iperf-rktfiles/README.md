# A simple iperf3 service

This directory contains scripts and files which can be used to create
a Rocket image which will provide a mongodb service.

## Source Files

##### mktree.sh
A script to compose the image file tree
##### manifest
The rocket image description file
##### run.sh
The task script. This file defines the set of operations that will be run inside the container

## Requirements

- actool
- rkt

## Image Structure

The image tree is composed under the current working directory in a
sub-directory named _image_.  It follows the AppContainer
specification for the location and format of the manifest file and the
location of the image contents.

    rootfs/
			/usr/local/bin/
				iperf3
			/lib64/
			  lib*.so
			/data/db/
			/dev/
			run.sh
### iperf3 daemon

The first file is the iperf3.  It is placed by convention into
`/usr/local/bin`.

## Building the File Tree

The _mktree.sh_ script will create the model tree for the image as
noted above.

    sh mktree.sh

The result is a directory named _image_ containing the image source
files.

## Building the Image

The image can be created merely using _tar_ but the AppContainer
project provides _actool_ to improve the process.  _actool_ will
validate the manifest before composing the image _.aci_ file.

    actool build image iperf3.aci

If you are rebuilding the image add _--overwrite_ after the _build_
argument.

Not that the version in the manifest must be incremented for each
buid.  The _rkt_ binary examines the version number in the image
manifest, and if it finds that they are the same, it assumes that it
can use the cached copy of the image.

== Running the Image

Right now this is a bit of a toy.  It runs as a single process inside
a container attached to a user shell

    sudo rkt run -volume=dev,kind=host,source=/dev iperf3.aci

## References

- https://github.com/appc/spec/blob/master/SPEC.md#image-layout[AppContainer Image Layout]
- https://github.com/appc/spec/blob/master/SPEC.md#image-manifest-schema[AppContainer manifest schema]
- https://github.com/appc/spec#building-acis[Using _actool_]
