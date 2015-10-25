ifndef BUILDTAGS
BUILDTAGS=predraft
endif
export BUILDTAGS

all:

	godep go build -tags $(BUILDTAGS) -o runtimeValidator

	$(MAKE) -C containerend
clean:
	rm runtimeValidator
	rm -rf  rootfs ubuntu.tar config.json runtime.json
	$(MAKE) -C containerend clean
