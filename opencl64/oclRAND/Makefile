# Builds mumax3 opencl kernels and create GO wrappers

SHELL = /bin/bash

OPENCLFILES := $(wildcard kernels_src/cl/*.cl)
WRAPPERFILES := $(notdir $(OPENCLFILES))
WRAPPERS := $(WRAPPERFILES:.cl=_wrapper.go)


.PHONY: all clean realclean


all:
	


clean: realclean
	rm -vf *_wrapper.go


realclean:
	rm -vf *_wrapper.go
