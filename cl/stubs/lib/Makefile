CC := gcc
CFLAGS := -fPIC -Wall -O4 -I../include
LDFLAGS := -shared

LIBPREFIX := 
LIBEXT := 
LINKS := 

ifeq ($(OS),Windows_NT)
	LIBEXT = lib
else
	LIBPREFIX = lib
	LIBEXT = so.1.2.0
	LINKS = ln -sf $(LIBPREFIX)OpenCL.${LIBEXT} $(LIBPREFIX)OpenCL.so
endif

.PHONY: all clean realclean

all:
	$(CC) $(LDFLAGS) $(CFLAGS) cl120.cc -o $(LIBPREFIX)OpenCL.${LIBEXT}
	$(LINKS)


clean:
	rm -f $(LIBPREFIX)OpenCL.*


realclean: clean

