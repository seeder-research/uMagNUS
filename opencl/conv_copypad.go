package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Copies src (larger) into dst (smaller).
// Used to extract demag field after convolution on padded m.
func copyUnPad(dst, src *data.Slice, dstsize, srcsize [3]int) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Argument(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	cfg := make3DConf(dstsize)

	event := k_copyunpad_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], cfg,
		[]*cl.Event{dst.GetEvent(0), src.GetEvent(0)})
	dst.SetEvent(0, event)
	src.SetEvent(0, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in copyunpad: %+v \n", err)
	}
}

// Copies src into dst, which is larger, and multiplies by vol*Bsat.
// The remainder of dst is not filled with zeros.
// Used to zero-pad magnetization before convolution and in the meanwhile multiply m by its length.
func copyPadMul(dst, src, vol *data.Slice, dstsize, srcsize [3]int, Msat MSlice) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	cfg := make3DConf(srcsize)

	event := k_copypadmul2_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z],
		Msat.DevPtr(0), Msat.Mul(0), vol.DevPtr(0), cfg,
		[]*cl.Event{dst.GetEvent(0), src.GetEvent(0), Msat.GetEvent(0), vol.GetEvent(0)})
	dst.SetEvent(0, event)
	src.SetEvent(0, event)
	Msat.SetEvent(0, event)
	vol.SetEvent(0, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in copypadmul: %+v \n", err)
	}
}
