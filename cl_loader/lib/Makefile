CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared

LIBPREFIX := 
LIBSUFFIX := 
LINKLIBS := 
ifeq ($(OS), WINDOWS_NT)
	LIBSUFFIX = lib
else
	LIBPREFIX = lib
	LIBSUFFIX = so.2.1.0
	LINKLIBS = ln -sf libumagnus.so.2.1.0 libumagnus.so.2.1 && ln -sf libumagnus.so.2.1 libumagnus.so.2 && ln -sf libumagnus.so.2 libumagnus.so && ln -sf libumagnus.so.2 libumagnus64.so
endif

.PHONY = all obj lib clean realclean

lib: obj
	$(CC) $(LDFLAGS) -o libumagnus.so.2.1.0 clCompiler.o
	$(CC) $(LDFLAGS) -o libumagnus64.so.2.1.0 clCompiler.o


all: lib
	$(LINKLIBS)


obj:
	$(MAKE) -C ../src


clean:
	rm -f *.o lib*.so lib*.so.*


realclean: clean
