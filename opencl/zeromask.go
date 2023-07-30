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
		event := k_zeromask_async(dst.DevPtr(c), unsafe.Pointer(mask), regions.Ptr, N, cfg, eql, q)

		dst.SetEvent(c, eventList[c])

		regions.InsertReadEvent(eventList[c])
		if Synchronous || Debug {
			if err := cl.WaitForEvents(eventList); err != nil {
				fmt.Printf("WaitForEvents failed in zeromask: %+v \n", err)
			}
		}
	}

	return
}
