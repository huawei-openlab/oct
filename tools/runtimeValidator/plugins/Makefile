
all:
	echo "Installing plugin: github.com/zenlinTechnofreak/ocitools..."
	set -e
	go get github.com/zenlinTechnofreak/ocitools
	go build github.com/zenlinTechnofreak/ocitools
	go build github.com/zenlinTechnofreak/ocitools/cmd/runtimetest
clean:
	rm -rf ocitools runtimetest 
	#rm -rf  rootfs ubuntu.tar config.json runtime.json
	#$(MAKE) -C containerend clean
