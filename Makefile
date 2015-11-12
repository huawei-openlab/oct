
all:
	go build -o ocitest
	$(MAKE) -C plugins
clean:
	go clean
	rm -rf ocitest
	$(MAKE) -C plugins clean
	$(MAKE) -C bundles clean
