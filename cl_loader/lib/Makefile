CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared

LIBPREFIX := 
LIBSUFFIX := 
LINKLIBS := 
LIBFILES := 
VER_MAJOR := 2
VER_MINOR := 4
VER_PATCH := 0
MINOR_VER := $(VER_MAJOR).$(VER_MINOR)
FULL_VER := $(MINOR_VER).$(VER_PATCH)
ifeq ($(OS),Windows_NT)
	LIBSUFFIX = dll
	LIBFILES = *.o *.dll
else
	LIBPREFIX = lib
	LIBSUFFIX = so.$(FULL_VER)
	LIBFILES = *.o lib*.so lib*.so*
	LINKLIBS = ln -sf libumagnus.so.$(FULL_VER) libumagnus.so.$(MINOR_VER) && ln -sf libumagnus.so.$(MINOR_VER) libumagnus.so.$(VER_MAJOR) && ln -sf libumagnus.so.$(VER_MAJOR) libumagnus.so && ln -sf libumagnus.so.$(VER_MAJOR) libumagnus64.so
endif

BUILD_TARGETS = all obj lib clean realclean


.PHONY: $(BUILD_TARGETS)

lib: obj
	$(CC) $(LDFLAGS) -o $(LIBPREFIX)umagnus.$(LIBSUFFIX) clCompiler.o
	$(CC) $(LDFLAGS) -o $(LIBPREFIX)umagnus64.$(LIBSUFFIX) clCompiler.o


all: lib
	$(LINKLIBS)


obj:
	$(MAKE) -C ../src


clean:
	rm -f $(LIBFILES)


realclean: clean
