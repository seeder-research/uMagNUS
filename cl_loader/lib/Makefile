CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared

.PHONY = all obj lib clean realclean

lib: obj
	$(CC) $(LDFLAGS) -o libumagnus.so.2.1.0 clCompiler.o


all: lib
	ln -s libumagnus.so.2.1.0 libumagnus.so.2.1
	ln -s libumagnus.so.2.1 libumagnus.so.2
	ln -s libumagnus.so.2 libumagnus.so


obj:
	$(MAKE) -C ../src


clean:
	rm -f *.o lib*.so lib*.so.*


realclean: clean
