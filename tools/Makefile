SUBDIR := octest #specsValidator

all:  
	@echo "Start to build." 
	@for subdir in $(SUBDIR); do\
		echo "making $$subdir";\
		$(MAKE) -C $$subdir;\
	done;
	@echo "Building finish."

.PHONY: clean

clean:
	@echo "Start to clean." 
	@for subdir in $(SUBDIR); do\
		$(MAKE) -C $$subdir clean;\
	done;
	@echo "Clean all!" 
