all:

	godep go build -o specsValidator

	$(MAKE) -C containerend
clean:
	rm specsValidator
	$(MAKE) -C containerend clean
