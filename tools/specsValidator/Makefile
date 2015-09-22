ifndef BUILDTAGS
BUILDTAGS=predraft
endif
export BUILDTAGS

all:

	godep go build -tags $(BUILDTAGS) -o specsValidator

	$(MAKE) -C containerend
clean:
	rm specsValidator
	rm -rf  rootfs ubuntu.tar config.json runtime.json
	$(MAKE) -C containerend clean
