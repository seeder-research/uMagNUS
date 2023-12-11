CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared
LIBPREFIX := 
LIBSUFFIX := 
LIBLINKS := 
LIB64LINKS := 
VER_MAJOR := 2
VER_MINOR := 4
VER_PATCH := 0
MINOR_VER := $(VER_MAJOR).$(VER_MINOR)
FULL_VER := $(MINOR_VER).$(VER_PATCH)
ifeq ($(OS),Windows_NT)
        LIBSUFFIX = dll
else
        LIBPREFIX = lib
        LIBSUFFIX = so.$(FULL_VER)
        LIBLINKS = ln -sf libumagnus.so.$(FULL_VER) libumagnus.so.$(MINOR_VER) && ln -sf libumagnus.so.$(MINOR_VER) libumagnus.so.$(VER_MAJOR) && ln -sf libumagnus.so.$(VER_MAJOR) libumagnus.so
        LIB64LINKS = ln -sf libumagnus64.so.$(FULL_VER) libumagnus64.so.$(MINOR_VER) && ln -sf libumagnus64.so.$(MINOR_VER) libumagnus64.so.$(VER_MAJOR) && ln -sf libumagnus64.so.$(VER_MAJOR) libumagnus64.so
endif

BUILD_TARGETS = all obj obj64 lib lib64 clean realclean

.PHONY: $(BUILD_TARGETS)

lib: obj
	$(CC) $(LDFLAGS) -o $(LIBPREFIX)umagnus.$(LIBSUFFIX) libumagnus.o
	$(LIBLINKS)


lib64: obj64
	$(CC) $(LDFLAGS) -o $(LIBPREFIX)umagnus64.$(LIBSUFFIX) libumagnus64.o
	$(LIB64LINKS)


all: lib lib64


obj: libumagnus.cc
	$(CC) $(CFLAGS) -o libumagnus.o -c libumagnus.cc


obj64: libumagnus64.cc
	$(CC) $(CFLAGS) -o libumagnus64.o -c libumagnus64.cc


clean:
	rm -f *.o lib*.so lib*.so.*


realclean: clean
	rm -f libumagnus.cc libumagnus64.cc
