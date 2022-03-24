# Use the default go compiler
GO_BUILDFLAGS=-compiler gc
# Or uncomment the line below to use the gccgo compiler, which may
# or may not be faster than gc and which may or may not compile...
# GO_BUILDFLAGS=-compiler gccgo -gccgoflags '-static-libgcc -O4 -Ofast -march=native'

CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'


all: 6g gccgo doc

install:
	go install -v $(GO_BUILDFLAGS)

6g:
	go install -v
	go vet *.go
	gofmt -w *.go

GCCGO=gccgo -gccgoflags '-static-libgcc -O3'

gccgo:
	go build -v -compiler $(GCCGO)

test: 6gtest gccgotest

6gtest:
	go test

gccgotest:
	go test -compiler $(GCCGO)

bench: 6gbench gccgobench

6gbench:
	go test -bench=.

gccgobench:
	go test -bench=. -compiler $(GCCGO)

clean:
	go clean

doc:
	godoc github.com/barnex/cuda5/cu > README
