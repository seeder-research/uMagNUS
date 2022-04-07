BUILDDIR="build"
STATIC="static"

# build the html pages in ${BUILDDIR}
.PHONY: html
html: doc uMagNUSlibs
	mkdir -p ${BUILDDIR}
	./doc -examples -builddir ${BUILDDIR}
	cp ${STATIC}/* build

.PHONY: doc
doc:
	go build -v

.PHONY: mumax3libs
uMagNUSlibs:
	$(MAKE) -C $(GOPATH)/src/github.com/seeder-research/uMagNUS uMagNUS

.PHONY: clean
clean:
	rm -rf build
	rm -f doc
