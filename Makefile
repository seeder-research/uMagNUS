
# Use the default go compiler
GO_BUILDFLAGS=-compiler gc
# Or uncomment the line below to use the gccgo compiler, which may
# or may not be faster than gc and which may or may not compile...
# GO_BUILDFLAGS=-compiler gccgo -gccgoflags '-static-libgcc -O4 -Ofast -march=native'

CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'


.PHONY: all cl-compiler clkernels clean realclean hooks go.mod

all: go.mod clkernels
	go install -v $(GO_BUILDFLAGS) github.com/seeder-research/uMagNUS/cmd/...

go.mod:
	go mod init

cl-compiler: go.mod
	go install -v github.com/seeder-research/uMagNUS/cmd/uMagNUS-clCompiler

clkernels:
	$(MAKE) -C ./opencl all
	$(MAKE) -C ./opencl64 all

hooks: .git/hooks/post-commit .git/hooks/pre-commit

.git/hooks/post-commit: post-commit
	ln -sf $(CURDIR)/$< $@

.git/hooks/pre-commit: pre-commit
	ln -sf $(CURDIR)/$< $@

clean:
	rm -frv $(GOPATH)/pkg/*/github.com/seeder-research/uMagNUS/*
	rm -frv $(GOPATH)/bin/mumax3* $(GOPATH)/bin/uMagNUS* go.mod
	$(MAKE) -C ./opencl clean
	$(MAKE) -C ./opencl64 clean

realclean: clean
	${MAKE} -C ./opencl realclean
	${MAKE} -C ./opencl64 realclean
