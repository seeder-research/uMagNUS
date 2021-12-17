# Builds mumax3 opencl kernels and create GO wrappers

SHELL = /bin/bash

.PHONY: all double opencl2go clean realclean


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
