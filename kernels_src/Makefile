# Builds uMagNUS opencl kernels and create GO wrappers

SHELL = /bin/bash

BUILD_TARGETS = all double opencl2go clean realclean


.PHONY: $(BUILD_TARGETS)


all: opencl2go
	./opencl2go


double: opencl2go
	./opencl2go -double


opencl2go: opencl2go.go
	go build $<


clean:
	rm -vf opencl2go


realclean:
	rm -vf opencl2go
