package opencl

import (
	"fmt"
	"unsafe"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl/cl"
)

// Sets vector dst to zero where mask != 0.
func ZeroMask(dst *data.Slice, mask LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)

	eventList := make([]*cl.Event, dst.NComp())
	for c := 0; c < dst.NComp(); c++ {
		eventList[c] = k_zeromask_async(dst.DevPtr(c), unsafe.Pointer(mask), regions.Ptr, N, cfg, [](*cl.Event){dst.GetEvent(c)})
		dst.SetEvent(c, eventList[c])
	}

	if err := cl.WaitForEvents(eventList); err != nil {
		fmt.Printf("WaitForEvents failed in zeromask: %+v \n", err)
	}
}
