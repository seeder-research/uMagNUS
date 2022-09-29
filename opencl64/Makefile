# Builds uMagNUS opencl kernels and create GO wrappers

SHELL = /bin/bash

OPENCLFILES := $(wildcard ../kernels_src/cl/*.cl)
WRAPPERFILES := $(notdir $(OPENCLFILES))
WRAPPERS := $(WRAPPERFILES:.cl=_wrapper.go)

OPENCL64FILES := $(wildcard ../kernels_src/cl64/*.cl)
WRAPPER64FILES := $(notdir $(OPENCL64FILES))
WRAPPERS64 := $(WRAPPER64FILES:.cl=_fp64_wrapper.go)

BUILD_TARGETS = all ocl2go wrappers fp64wrappers clean realclean


.PHONY: $(BUILD_TARGETS)


all: ocl2go wrappers fp64wrappers opencl2go
	mv RNGthreefry*_wrapper.go oclRAND/.
	mv RNGxorwow*_wrapper.go oclRAND/.
	sed -i "s/package opencl64/package oclRAND/g" oclRAND/RNGthreefry*wrapper.go
	sed -i "s/package opencl64/package oclRAND/g" oclRAND/RNGxorwow*wrapper.go
	go install -v


wrappers: $(WRAPPERS)


fp64wrappers: $(WRAPPERS64)


ocl2go:
	$(MAKE) -C ../ocl2go all


opencl2go: ocl2go
	$(MAKE) -C ../kernels_src clean
	$(MAKE) -C ../kernels_src double


%_wrapper.go: ../kernels_src/cl/%.cl
	../ocl2go/ocl2go -outdir="." -lib64 $< > /dev/null
	gofmt -w -s -l ./$@ > /dev/null


%_fp64_wrapper.go: ../kernels_src/cl64/%.cl
	../ocl2go/ocl2go -outdir="." -lib64 -double $< > /dev/null
	gofmt -w -s -l ./$@ > /dev/null


clean:
	$(MAKE) -C ../kernels_src clean
	$(MAKE) -C ./oclRAND clean
	rm -vf *_wrapper.go ../kernels64/program_wrapper.go


realclean:
	$(MAKE) -C ../kernels_src realclean
	rm -vf *_wrapper.go ../kernels64/program_wrapper.go
