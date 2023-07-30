package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// dst += LUT[region], for vectors. Used to add terms to excitation.
func RegionAddV(dst *data.Slice, lut LUTPtrs, regions *Bytes, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(dst.NComp() == 3)
	N := dst.Len()
	cfg := make1DConf(N)

	event := k_regionaddv_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		lut[X], lut[Y], lut[Z], regions.Ptr, N, cfg, ewl, q)

	dst.SetEvent(X, event)
	dst.SetEvent(Y, event)
	dst.SetEvent(Z, event)

	regions.InsertReadEvent(event)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in regionaddv failed: %+v \n", err)
		}
		regions.RemoveReadEvent(event)
	}

	return
}

// dst += LUT[region], for scalar. Used to add terms to scalar excitation.
func RegionAddS(dst *data.Slice, lut LUTPtr, regions *Bytes, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(dst.NComp() == 1)
	N := dst.Len()
	cfg := make1DConf(N)

	event := k_regionadds_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		ewl, q)

	dst.SetEvent(0, event)
	regions.InsertReadEvent(event)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in regionadds failed: %+v \n", err)
		}
		regions.RemoveReadEvent(event)
	}

	return
}

// decode the regions+LUT pair into an uncompressed array
func RegionDecode(dst *data.Slice, lut LUTPtr, regions *Bytes, q *cl.CommandQueue, ewl []*cl.Event) {
	N := dst.Len()
	cfg := make1DConf(N)

	event := k_regiondecode_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		ewl, q)

	dst.SetEvent(0, event)
	regions.InsertReadEvent(event)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in regiondecode failed: %+v \n", err)
		}
		regions.RemoveReadEvent(event)
	}

	return
}

// select the part of src within the specified region, set 0's everywhere else.
func RegionSelect(dst, src *data.Slice, regions *Bytes, region byte, q []*cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(dst.NComp() == src.NComp())
	util.Argument(dst.NComp() == len(q))
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		event := k_regionselect_async(dst.DevPtr(c), src.DevPtr(c), regions.Ptr, region, N, cfg,
			ewl, q[c])

		dst.SetEvent(c, event)
		src.InsertReadEvent(c, event)
		regions.InsertReadEvent(event)

		if Synchronous || Debug {
			if err := cl.WaitForEvents(event); err != nil {
				fmt.Printf("WaitForEvents in regionselect failed: %+v \n", err)
			}
		}
	}

	return
}
