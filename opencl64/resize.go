package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

// Select and resize one layer for interactive output
func Resize(dst, src *data.Slice, layer int) {
	dstsize := dst.Size()
	srcsize := src.Size()
	util.Assert(dstsize[Z] == 1)
	util.Assert(dst.NComp() == 1 && src.NComp() == 1)

	scalex := srcsize[X] / dstsize[X]
	scaley := srcsize[Y] / dstsize[Y]
	util.Assert(scalex > 0 && scaley > 0)

	cfg := make3DConf(dstsize)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent(X)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = dst.GetEvent(Y)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = dst.GetEvent(Z)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = src.GetEvent(X)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = src.GetEvent(Y)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = src.GetEvent(Z)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_resize_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], layer, scalex, scaley, cfg,
		eventsList)

	dst.SetEvent(X, event)
	dst.SetEvent(Y, event)
	dst.SetEvent(Z, event)

	glist := []GSlice{src}
	InsertEventIntoGSlices(event, glist)

	// Synchronize for resize
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in resize: %+v \n", err)
		WaitAndUpdateDataSliceEvents(event, glist, false)
	}
}
