package opencl

import (
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Sets vector dst to zero where mask != 0.
func ZeroMask(dst *data.Slice, mask LUTPtr, regions *Bytes, queue []*cl.CommandQueue, events []*cl.Event) {
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		k_zeromask_async(dst.DevPtr(c), unsafe.Pointer(mask),
			regions.Ptr, N,
			cfg, queue[c], events)
	}
}
