# Builds mumax3 opencl kernels and create GO wrappers

SHELL = /bin/bash

OPENCLFILES := $(wildcard ../kernels_src/cl/*.cl)
WRAPPERFILES := $(notdir $(OPENCLFILES))
WRAPPERS := $(WRAPPERFILES:.cl=_wrapper.go)


.PHONY: all ocl2go wrappers clean realclean


all: ocl2go wrappers opencl2go
	mv RNGmtgp*_wrapper.go oclRAND/.
	mv RNGthreefry*_wrapper.go oclRAND/.
	mv RNGxorwow*_wrapper.go oclRAND/.
	sed -i "s/package opencl/package oclRAND/g" oclRAND/RNGmtgp*wrapper.go
	sed -i "s/package opencl/package oclRAND/g" oclRAND/RNGthreefry*wrapper.go
	sed -i "s/package opencl/package oclRAND/g" oclRAND/RNGxorwow*wrapper.go
	go install -v


wrappers: $(WRAPPERS)


ocl2go: ocl2go.go
	go build $<


opencl2go: ocl2go
	$(MAKE) -C ../kernels_src clean
	$(MAKE) -C ../kernels_src


%_wrapper.go: ../kernels_src/cl/%.cl ocl2go
	./ocl2go $< > /dev/null
	gofmt -w -s -l ../kernels_src/cl/$@ > /dev/null
	mv ../kernels_src/cl/$@ ./$@


clean:
	$(MAKE) -C ../kernels_src clean
	$(MAKE) -C ./oclRAND clean
	rm -vf *_wrapper.go ocl2go ../kernels/program_wrapper.go


realclean:
	$(MAKE) -C ../kernels_src realclean
	rm -vf *_wrapper.go ocl2go ../kernels/program_wrapper.go
