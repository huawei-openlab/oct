ifndef BUILDTAGS
BUILDTAGS=v0.1.1
endif
export BUILDTAGS

all:

	godep go build -tags $(BUILDTAGS) -o specsValidator

	$(MAKE) -C containerend
clean:
	rm specsValidator
	$(MAKE) -C containerend clean
