# Builds uMagNUS opencl kernels and create GO wrappers

SHELL = /bin/bash

OPENCLFILES := $(wildcard ../kernels_src/cl/*.cl)
WRAPPERFILES := $(notdir $(OPENCLFILES))
WRAPPERS := $(WRAPPERFILES:.cl=_wrapper.go)

BUILD_TARGETS = all ocl2go wrappers clean realclean


.PHONY: $(BUILD_TARGETS)


all: opencl2go wrappers
	mv RNGthreefry*_wrapper.go oclRAND/.
	mv RNGxorwow*_wrapper.go oclRAND/.
	sed -i "s/package opencl/package oclRAND/g" oclRAND/RNGthreefry*wrapper.go
	sed -i "s/package opencl/package oclRAND/g" oclRAND/RNGxorwow*wrapper.go
	go install -v


wrappers: $(WRAPPERS)


ocl2go:
	$(MAKE) -C ../ocl2go all


opencl2go: ocl2go
	$(MAKE) -C ../kernels_src clean
	$(MAKE) -C ../kernels_src


%_wrapper.go: ../kernels_src/cl/%.cl ocl2go
	../ocl2go/ocl2go -outdir="." $< > /dev/null
	gofmt -w -s -l ./$@ > /dev/null


clean:
	$(MAKE) -C ../kernels_src clean
	$(MAKE) -C ./oclRAND clean
	rm -vf *_wrapper.go ../kernels/program_wrapper.go


realclean:
	$(MAKE) -C ../kernels_src realclean
	rm -vf *_wrapper.go ../kernels/program_wrapper.go
