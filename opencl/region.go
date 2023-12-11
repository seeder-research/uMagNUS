package opencl

import (
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// dst += LUT[region], for vectors. Used to add terms to excitation.
func RegionAddV(dst *data.Slice, lut LUTPtrs, regions *Bytes, queue *cl.CommandQueue, events []*cl.Event) {
	util.Argument(dst.NComp() == 3)
	N := dst.Len()
	cfg := make1DConf(N)

	k_regionaddv_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		lut[X], lut[Y], lut[Z], regions.Ptr, N, cfg, queue, events)
}

// dst += LUT[region], for scalar. Used to add terms to scalar excitation.
func RegionAddS(dst *data.Slice, lut LUTPtr, regions *Bytes, queue *cl.CommandQueue, events []*cl.Event) {
	util.Argument(dst.NComp() == 1)
	N := dst.Len()
	cfg := make1DConf(N)

	k_regionadds_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		queue, events)
}

// decode the regions+LUT pair into an uncompressed array
func RegionDecode(dst *data.Slice, lut LUTPtr, regions *Bytes, queue *cl.CommandQueue, events []*cl.Event) {
	N := dst.Len()
	cfg := make1DConf(N)

	k_regiondecode_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		queue, events)
}

// select the part of src within the specified region, set 0's everywhere else.
func RegionSelect(dst, src *data.Slice, regions *Bytes, region byte, queue []*cl.CommandQueue, events []*cl.Event) {
	util.Argument(dst.NComp() == src.NComp())
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		k_regionselect_async(dst.DevPtr(c), src.DevPtr(c), regions.Ptr, region, N, cfg,
			queue[c], events)
	}
}
