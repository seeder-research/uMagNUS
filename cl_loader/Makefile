CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared

BUILD_TARGETS = all clean realclean


.PHONY: $(BUILD_TARGETS)


all:
	$(MAKE) -C ./lib all


clean:
	$(MAKE) -C ./lib clean


realclean:
	$(MAKE) -C ./lib realclean
