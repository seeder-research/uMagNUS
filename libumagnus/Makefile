CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared

.PHONY = all obj lib clean realclean

lib: obj
	$(CC) $(LDFLAGS) -o libumagnus.so.2.1.0 libumagnus.o


all: lib
	ln -sf libumagnus.so.2.1.0 libumagnus.so.2.1
	ln -sf libumagnus.so.2.1 libumagnus.so.2
	ln -sf libumagnus.so.2 libumagnus.so


obj: libumagnus.cc
	$(CC) $(CFLAGS) -o libumagnus.o -c libumagnus.cc


clean:
	rm -f *.o lib*.so lib*.so.*


realclean: clean
	rm -f libumagnus.cc
