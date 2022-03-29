package cl

// This file provides CGO flags to find OpecnCL libraries and headers.

//#cgo darwin LDFLAGS: -framework OpenCL
//#cgo !darwin LDFLAGS: -lOpenCL
//
////default location:
//#cgo LDFLAGS:-L./stubs/lib
//#cgo CFLAGS: -I./stubs/include
//
import "C"
