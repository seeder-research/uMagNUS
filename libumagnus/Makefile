CC := gcc
CFLAGS := -fPIC -Wall -O4
LDFLAGS := -shared

.PHONY = all obj obj64 lib lib64 clean realclean

lib: obj
	$(CC) $(LDFLAGS) -o libumagnus.so.2.1.0 libumagnus.o
	ln -sf libumagnus.so.2.1.0 libumagnus.so.2.1
	ln -sf libumagnus.so.2.1 libumagnus.so.2
	ln -sf libumagnus.so.2 libumagnus.so


lib64: obj64
	$(CC) $(LDFLAGS) -o libumagnus64.so.2.1.0 libumagnus64.o
	ln -sf libumagnus64.so.2.1.0 libumagnus64.so.2.1
	ln -sf libumagnus64.so.2.1 libumagnus64.so.2
	ln -sf libumagnus64.so.2 libumagnus64.so


all: lib lib64


obj: libumagnus.cc
	$(CC) $(CFLAGS) -o libumagnus.o -c libumagnus.cc


obj64: libumagnus64.cc
	$(CC) $(CFLAGS) -o libumagnus64.o -c libumagnus64.cc


clean:
	rm -f *.o lib*.so lib*.so.*


realclean: clean
	rm -f libumagnus.cc libumagnus64.cc
