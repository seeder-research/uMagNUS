package opencl64

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

// Copies src (larger) into dst (smaller).
// Used to extract demag field after convolution on padded m.
func copyUnPad(dst, src *data.Slice, dstsize, srcsize [3]int) {
	dstPtr := (unsafe.Pointer)(nil)
	srcPtr := (unsafe.Pointer)(nil)
	eventList := []*cl.Event{}

	if dst != nil {
		dstPtr = dst.DevPtr(0)
		eventList = append(eventList, dst.GetEvent(0))
	} else {
		panic("ERROR (copyUnPad): dst pointer cannot be nil")
	}
	if src != nil {
		srcPtr = src.DevPtr(0)
		eventList = append(eventList, src.GetEvent(0))
	} else {
		panic("ERROR (copyUnPad): dst pointer cannot be nil")
	}

	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Argument(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	cfg := make3DConf(dstsize)

	event := k_copyunpad_async(dstPtr, dstsize[X], dstsize[Y], dstsize[Z],
		srcPtr, srcsize[X], srcsize[Y], srcsize[Z], cfg,
		eventList)

	dst.SetEvent(0, event)
	src.SetEvent(0, event)

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in copyunpad: %+v \n", err)
	}
}

// Copies src into dst, which is larger, and multiplies by vol*Bsat.
// The remainder of dst is not filled with zeros.
// Used to zero-pad magnetization before convolution and in the meanwhile multiply m by its length.
func copyPadMul(dst, src, vol *data.Slice, dstsize, srcsize [3]int, Msat MSlice) {
	dstPtr := (unsafe.Pointer)(nil)
	srcPtr := (unsafe.Pointer)(nil)
	volPtr := (unsafe.Pointer)(nil)
	MsPtr := (unsafe.Pointer)(nil)
	eventList := []*cl.Event{}

	if dst != nil {
		dstPtr = dst.DevPtr(0)
		eventList = append(eventList, dst.GetEvent(0))
	} else {
		panic("ERROR (copyPadMul): dst pointer cannot be nil")
	}
	if src != nil {
		srcPtr = src.DevPtr(0)
		eventList = append(eventList, src.GetEvent(0))
	} else {
		panic("ERROR (copyPadMul): src pointer cannot be nil")
	}
	if vol != nil {
		volPtr = vol.DevPtr(0)
		eventList = append(eventList, vol.GetEvent(0))
	}
	if Msat.GetSlicePtr(0) != nil {
		MsPtr = Msat.DevPtr(0)
		eventList = append(eventList, Msat.GetEvent(0))
	}

	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	cfg := make3DConf(srcsize)

	event := k_copypadmul2_async(dstPtr, dstsize[X], dstsize[Y], dstsize[Z],
		srcPtr, srcsize[X], srcsize[Y], srcsize[Z],
		MsPtr, Msat.Mul(0), volPtr, cfg,
		eventList)

	dst.SetEvent(0, event)
	src.SetEvent(0, event)
	if Msat.GetSlicePtr(0) != nil {
		Msat.SetEvent(0, event)
	}
	if vol != nil {
		vol.SetEvent(0, event)
	}

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in copypadmul: %+v \n", err)
	}
}
