package loader

// This file provides CGO flags to find uMagNUS libraries and headers.

//#cgo darwin LDFLAGS: -framework umagnus
//#cgo !darwin LDFLAGS: -lumagnus
//
////default location:
//#cgo LDFLAGS:-L${SRCDIR}/cl_loader/lib
//#cgo LDFLAGS:-L${SRCDIR}/cl_loader/lib64
//#cgo CFLAGS: -I${SRCDIR}/cl_loader/include
//
import "C"
