package opencl

import (
	"fmt"
	"unsafe"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl/cl"
	"github.com/seeder-research/uMagNUS/util"
)

// dst += LUT[region], for vectors. Used to add terms to excitation.
func RegionAddV(dst *data.Slice, lut LUTPtrs, regions *Bytes) {
	util.Argument(dst.NComp() == 3)
	N := dst.Len()
	cfg := make1DConf(N)
	event := k_regionaddv_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		lut[X], lut[Y], lut[Z], regions.Ptr, N, cfg,
		[](*cl.Event){dst.GetEvent(X), dst.GetEvent(Y), dst.GetEvent(Z)})
	dst.SetEvent(X, event)
	dst.SetEvent(Y, event)
	dst.SetEvent(Z, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents in regionaddv failed: %+v \n", err)
	}
}

// dst += LUT[region], for scalar. Used to add terms to scalar excitation.
func RegionAddS(dst *data.Slice, lut LUTPtr, regions *Bytes) {
	util.Argument(dst.NComp() == 1)
	N := dst.Len()
	cfg := make1DConf(N)
	event := k_regionadds_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		[](*cl.Event){dst.GetEvent(0)})
	dst.SetEvent(0, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents in regionadds failed: %+v \n", err)
	}
}

// decode the regions+LUT pair into an uncompressed array
func RegionDecode(dst *data.Slice, lut LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)
	event := k_regiondecode_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		[](*cl.Event){dst.GetEvent(0)})
	dst.SetEvent(0, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents in regiondecode failed: %+v \n", err)
	}
}

// select the part of src within the specified region, set 0's everywhere else.
func RegionSelect(dst, src *data.Slice, regions *Bytes, region byte) {
	util.Argument(dst.NComp() == src.NComp())
	N := dst.Len()
	cfg := make1DConf(N)

	eventList := make([]*cl.Event, dst.NComp())
	for c := 0; c < dst.NComp(); c++ {
		eventList[c] = k_regionselect_async(dst.DevPtr(c), src.DevPtr(c), regions.Ptr, region, N, cfg,
			[](*cl.Event){dst.GetEvent(c), src.GetEvent(c)})
		dst.SetEvent(c, eventList[c])
		src.SetEvent(c, eventList[c])
	}
	err := cl.WaitForEvents(eventList)
	if err != nil {
		fmt.Printf("WaitForEvents in regionselect failed: %+v \n", err)
	}
}
