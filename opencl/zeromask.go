package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Sets vector dst to zero where mask != 0.
func ZeroMask(dst *data.Slice, mask LUTPtr, regions *Bytes, q []*cl.CommandQueue, ewl []*cl.Event) {
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		// Launch kernel
		event := k_zeromask_async(dst.DevPtr(c), unsafe.Pointer(mask), regions.Ptr, N, cfg, ewl, q[c])

		if Debug {
			if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
				fmt.Printf("WaitForEvents failed in zeromask: %+v \n", err)
			}
		}
	}

	return
}
