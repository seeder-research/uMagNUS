# Use the default go compiler
GO_BUILDFLAGS=-compiler gc
# Or uncomment the line below to use the gccgo compiler, which may
# or may not be faster than gc and which may or may not compile...
# GO_BUILDFLAGS=-compiler gccgo -gccgoflags '-static-libgcc -O4 -Ofast -march=native'

SRCDIR := $(shell pwd)

GOPATH ?= $(SRCDIR)/gopath

BUILDPATH := $(GOPATH)

CLEANLIBFILES := 
ifeq ($(OS), Windows_NT)
	CLEANLIBFILES = rm -frv $(BUILDPATH)/bin/*.dll
endif

LIBUMAGNUS32SRC := $(PWD)/kernels_src/Kernels/kernels32.h
LIBUMAGNUS64SRC := $(PWD)/kernels_src/Kernels/kernels64.h

CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'

DIR_TARGET = $(BUILDPATH)

BUILD_TARGETS = all base mod cl-binds cl-compiler clkernels clean data data64 draw draw64 dump dump64 engine engine64 freetype gui realclean hooks httpfs mag mag64 oommf oommf64 script script64 timer uMagNUS uMagNUS64 util loader loader64 kernloader kernloader64 libumagnus libumagnus64 libs


.PHONY: $(BUILD_TARGETS)


.EXPORT_ALL_VARIABLES:
	GOPATH = $(BUILDPATH)


all: base libs


$(DIR_TARGET):
	mkdir -p $(BUILDPATH)


base: mod cl-compiler kernloader kernloader64 uMagNUS uMagNUS64
	go install -v $(GO_BUILDFLAGS) github.com/seeder-research/uMagNUS/cmd/...


hooks: .git/hooks/post-commit .git/hooks/pre-commit


.git/hooks/post-commit: post-commit
	ln -sf $(CURDIR)/$< $@


.git/hooks/pre-commit: pre-commit
	ln -sf $(CURDIR)/$< $@


mod: $(DIR_TARGET)
	go mod init github.com/seeder-research/uMagNUS


cl-binds: go.mod
	$(MAKE) -C ./cl install


clkernels: go.mod
	$(MAKE) -C ./opencl all


clkernels64: go.mod
	$(MAKE) -C ./opencl64 all


freetype: go.mod
	go install -v $(GO_BUILDFLAGS) github.com/seeder-research/uMagNUS/freetype/raster


gui: go.mod
	$(MAKE) -C ./gui all


httpfs: go.mod
	$(MAKE) -C ./httpfs all


timer: go.mod
	$(MAKE) -C ./timer all


util: go.mod
	$(MAKE) -C ./util all


ocl2go: go.mod
	$(MAKE) -C ./ocl2go all


cl-compiler: cl-binds
	$(MAKE) -C ./cmd/uMagNUS-clCompiler all


loader: cl-binds
	$(MAKE) -C ./cl_loader all
	$(MAKE) -C ./loader all


loader64: cl-binds
	$(MAKE) -C ./cl_loader all
	$(MAKE) -C ./loader64 all


loaders: loader loader64


kernloader: loader
	$(MAKE) -C ./cmd/uMagNUS-kernelLoader all


kernloader64: loader64
	$(MAKE) -C ./cmd/uMagNUS-kernelLoader64 all


libumagnus: cl-compiler
	rm -f ./libumagnus/*.cc
	$(BUILDPATH)/bin/uMagNUS-clCompiler -args="-cl-opt-disable -cl-mad-enable -cl-finite-math-only -cl-single-precision-constant -cl-fp32-correctly-rounded-divide-sqrt -cl-kernel-arg-info" -std="CL1.2" -iopts="-I$(PWD)/kernels_src" -dump $(LIBUMAGNUS32SRC) >> libumagnus/libumagnus.cc
	$(MAKE) -C ./libumagnus lib


libumagnus64: cl-compiler
	rm -f ./libumagnus/*.cc
	$(BUILDPATH)/bin/uMagNUS-clCompiler -args="-cl-opt-disable -cl-mad-enable -cl-finite-math-only -cl-fp32-correctly-rounded-divide-sqrt -cl-kernel-arg-info -D__REAL_IS_DOUBLE__" -std="CL1.2" -iopts="-I$(PWD)/kernels_src" -dump $(LIBUMAGNUS64SRC) >> libumagnus/libumagnus64.cc
	$(MAKE) -C ./libumagnus lib64


libs: libumagnus libumagnus64


data: cl-binds util
	$(MAKE) -C ./data all


data64: cl-binds util
	$(MAKE) -C ./data64 all


script: data
	$(MAKE) -C ./script all


script64: data64
	$(MAKE) -C ./script64 all


draw: data freetype util
	$(MAKE) -C ./draw all


draw64: data64 freetype util
	$(MAKE) -C ./draw64 all


dump: data util
	$(MAKE) -C ./dump all


dump64: data64 util
	$(MAKE) -C ./dump64 all


oommf: data util
	$(MAKE) -C ./oommf all


oommf64: data64 util
	$(MAKE) -C ./oommf64 all


mag: oommf timer
	$(MAKE) -C ./mag all


mag64: oommf64 timer
	$(MAKE) -C ./mag64 all


engine: clkernels gui httpfs mag script loader
	$(MAKE) -C ./engine all


engine64: clkernels64 gui httpfs mag64 script64 loader64
	$(MAKE) -C ./engine64 all


uMagNUS: engine
	$(MAKE) -C ./cmd/uMagNUS all


uMagNUS64: engine64
	$(MAKE) -C ./cmd/uMagNUS64 all


clean:
	rm -frv $(BUILDPATH)/pkg/*/github.com/seeder-research/uMagNUS/*
	rm -frv $(BUILDPATH)/bin/mumax3* $(BUILDPATH)/bin/uMagNUS* go.mod
	$(CLEANLIBFILES)
	$(MAKE) -C ./cl_loader clean
	$(MAKE) -C ./opencl clean
	$(MAKE) -C ./opencl64 clean
	$(MAKE) -C ./kernels_src/Kernels clean
	$(MAKE) -C ./libumagnus clean
	$(MAKE) -C ./ocl2go realclean
	$(MAKE) -C ./cl/stubs clean


realclean: clean
	rm -frv ./gopath
	${MAKE} -C ./cl_loader realclean
	${MAKE} -C ./opencl realclean
	${MAKE} -C ./opencl64 realclean
	$(MAKE) -C ./kernels_src/Kernels realclean
	$(MAKE) -C ./libumagnus realclean
	$(MAKE) -C ./ocl2go realclean
	$(MAKE) -C ./cl/stubs realclean
