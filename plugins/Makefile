
all:
	echo ${GOPATH}
	echo "Installing plugin: github.com/zenlinTechnofreak/ocitools..."
	set -e
	go get github.com/zenlinTechnofreak/ocitools
	go build github.com/zenlinTechnofreak/ocitools
	go build github.com/zenlinTechnofreak/ocitools/cmd/runtimetest
	go get github.com/huawei-openlab/oci2aci
	go build github.com/huawei-openlab/oci2aci
clean:
	go clean
	#rm -rf ${GOPATH}/src/github.com/zenlinTechnofreak/ocitools
	rm -rf ocitools runtimetest oci2aci
	rm -rf config.json* runtime.json*
	#$(MAKE) -C containerend clean
