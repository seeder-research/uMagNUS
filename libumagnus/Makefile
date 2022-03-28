CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared
LIBPREFIX := 
LIBSUFFIX := 
LIBLINKS := 
LIB64LINKS := 
ifeq ($(OS),Windows_NT)
        LIBSUFFIX = dll
else
        LIBPREFIX = lib
        LIBSUFFIX = so.2.1.0
        LIBLINKS = ln -sf libumagnus.so.2.1.0 libumagnus.so.2.1 && ln -sf libumagnus.so.2.1 libumagnus.so.2 && ln -sf libumagnus.so.2 libumagnus.so
        LIB64LINKS = ln -sf libumagnus64.so.2.1.0 libumagnus64.so.2.1 && ln -sf libumagnus64.so.2.1 libumagnus64.so.2 && ln -sf libumagnus64.so.2 libumagnus64.so
endif

.PHONY = all obj obj64 lib lib64 clean realclean

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
