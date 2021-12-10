# Builds mumax3 opencl kernels and create GO wrappers

SHELL = /bin/bash

OPENCLFILES := $(wildcard ../kernels_src/cl/*.cl)
WRAPPERFILES := $(notdir $(OPENCLFILES))
WRAPPERS := $(WRAPPERFILES:.cl=_wrapper.go)


.PHONY: all ocl2go64 wrappers clean realclean


all: ocl2go64 wrappers opencl2go
	mv RNGmrg32k3a*_wrapper.go oclRAND/.
	mv RNGmtgp*_wrapper.go oclRAND/.
	mv RNGthreefry*_wrapper.go oclRAND/.
	mv RNGxorwow*_wrapper.go oclRAND/.
	sed -i "s/package opencl64/package oclRAND/g" oclRAND/RNGmrg32k3a*wrapper.go
	sed -i "s/package opencl64/package oclRAND/g" oclRAND/RNGmtgp*wrapper.go
	sed -i "s/package opencl64/package oclRAND/g" oclRAND/RNGthreefry*wrapper.go
	sed -i "s/package opencl64/package oclRAND/g" oclRAND/RNGxorwow*wrapper.go
	go install -v


wrappers: $(WRAPPERS)


ocl2go64: ocl2go64.go
	go build $<


opencl2go: ocl2go64
	$(MAKE) -C ../kernels_src
	mv ../kernels_src/opencl_kernels_wrapper.go ./
	sed -i "s/package opencl/package opencl64/" opencl_kernels_wrapper.go


%_wrapper.go: ../kernels_src/cl/%.cl ocl2go64
	./ocl2go64 $< > /dev/null
	gofmt -w -s -l ../kernels_src/cl/$@ > /dev/null
	mv ../kernels_src/cl/$@ ./$@


clean:
	$(MAKE) -C ../kernels_src clean
	$(MAKE) -C ./oclRAND clean
	rm -vf *_wrapper.go ocl2go64 ../kernels/program_wrapper.go


realclean:
	$(MAKE) -C ../kernels_src realclean
	rm -vf *_wrapper.go ocl2go64 ../kernels/program_wrapper.go
