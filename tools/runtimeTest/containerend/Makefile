SOURCES:=$(wildcard *.go)
TARGETLIST:=$(patsubst %.go,%,$(SOURCES))
all:
	$(foreach TARGET,${TARGETLIST},godep go build $(TARGET).go;)
.PHONY:
	clean
clean:
	$(foreach TARGET,${TARGETLIST},rm -f $(TARGET);)



