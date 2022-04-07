BUILDDIR="build"
STATIC="static"

# build the html pages in ${BUILDDIR}
.PHONY: doc uMagNUSlibs html clean

html: uMagNUSlibs doc
	mkdir -p ${BUILDDIR}
	./doc -examples -builddir ${BUILDDIR}
	cp ${STATIC}/* build


doc:
	go build -v


uMagNUSlibs:
	cd $(GOPATH)/src/github.com/seeder-research/uMagNUS && $(MAKE) cl-binds clkernels engine uMagNUS libumagnus clkernels64 engine64 uMagNUS64 libumagnus64 && go install -v github.com/seeder-research/uMagNUS/cmd/...


clean:
	rm -rf build
	rm -f doc
