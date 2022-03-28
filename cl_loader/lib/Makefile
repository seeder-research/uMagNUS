CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared

LIBPREFIX := 
LIBSUFFIX := 
LINKLIBS := 
LIBFILES := 
ifeq ($(OS),Windows_NT)
	LIBSUFFIX = dll
	LIBFILES = *.o *.dll
else
	LIBPREFIX = lib
	LIBSUFFIX = so.2.1.0
	LIBFILES = *.o lib*.so lib*.so*
	LINKLIBS = ln -sf libumagnus.so.2.1.0 libumagnus.so.2.1 && ln -sf libumagnus.so.2.1 libumagnus.so.2 && ln -sf libumagnus.so.2 libumagnus.so && ln -sf libumagnus.so.2 libumagnus64.so
endif

.PHONY = all obj lib clean realclean

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
