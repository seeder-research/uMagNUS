# Builds mumax3 opencl kernels and create GO wrappers

SHELL = /bin/bash

OPENCLFILES := $(wildcard kernels_src/cl/*.cl)
WRAPPERFILES := $(notdir $(OPENCLFILES))
WRAPPERS := $(WRAPPERFILES:.cl=_wrapper.go)


.PHONY: all ocl2go wrappers clean realclean


all: ocl2go wrappers opencl2go
	mv RNGmtgp*_wrapper.go RNGmtgp/.
	sed -i "s/package opencl/package RNGmtgp/g" RNGmtgp/RNGmtgp*.go
	go install -v


wrappers: $(WRAPPERS)


ocl2go: ocl2go.go
	go build $<


opencl2go: ocl2go
	cd ./kernels_src && $(MAKE)


%_wrapper.go: kernels_src/cl/%.cl ocl2go
	./ocl2go $< > /dev/null
	gofmt -w -s -l ./kernels_src/cl/$@ > /dev/null
	mv ./kernels_src/cl/$@ ./$@


clean:
	cd ./kernels_src && $(MAKE) clean
	cd ./RNGmtgp && $(MAKE) clean
	rm -vf *_wrapper.go ocl2go ./kernels/program_wrapper.go


realclean:
	cd ./kernels_src && $(MAKE) realclean
	rm -vf *_wrapper.go ocl2go ./kernels/program_wrapper.go
