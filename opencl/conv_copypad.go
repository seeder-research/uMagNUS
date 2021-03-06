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

	eventList := []*cl.Event{}
	tmpEvt := dst.GetEvent(0)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = src.GetEvent(0)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_copyunpad_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], cfg,
		eventList)

	dst.SetEvent(0, event)
	src.SetEvent(0, event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in copyunpad: %+v \n", err)
		}
	}
}

// Copies src into dst, which is larger, and multiplies by vol*Bsat.
// The remainder of dst is not filled with zeros.
// Used to zero-pad magnetization before convolution and in the meanwhile multiply m by its length.
func copyPadMul(dst, src, vol *data.Slice, dstsize, srcsize [3]int, Msat MSlice) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	cfg := make3DConf(srcsize)

	eventList := []*cl.Event{}
	tmpEvent := dst.GetEvent(0)
	if tmpEvent != nil {
		eventList = append(eventList, tmpEvent)
	}
	tmpEvent = src.GetEvent(0)
	if tmpEvent != nil {
		eventList = append(eventList, tmpEvent)
	}
	tmpEvent = vol.GetEvent(0)
	if tmpEvent != nil {
		eventList = append(eventList, tmpEvent)
	}
	if Msat.GetSlicePtr() != nil {
		tmpEvent = Msat.GetEvent(0)
		if tmpEvent != nil {
			eventList = append(eventList, tmpEvent)
		}
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_copypadmul2_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z],
		Msat.DevPtr(0), Msat.Mul(0), vol.DevPtr(0), cfg,
		eventList)

	dst.SetEvent(0, event)
	src.SetEvent(0, event)
	vol.SetEvent(0, event)
	if Msat.GetSlicePtr() != nil {
		Msat.SetEvent(0, event)
	}

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in copypadmul: %+v \n", err)
		}
	}
}
