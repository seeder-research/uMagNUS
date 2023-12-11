package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Select and resize one layer for interactive output
func Resize(dst, src *data.Slice, layer int, queue *cl.CommandQueue, events []*cl.Event) {
	dstsize := dst.Size()
	srcsize := src.Size()
	util.Assert(dstsize[Z] == 1)
	util.Assert(dst.NComp() == 1 && src.NComp() == 1)

	scalex := srcsize[X] / dstsize[X]
	scaley := srcsize[Y] / dstsize[Y]
	util.Assert(scalex > 0 && scaley > 0)

	cfg := make3DConf(dstsize)

	event := k_resize_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], layer, scalex, scaley, cfg,
		queue, events)

	// Synchronize for resize (?need to?)
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in resize: %+v \n", err)
	}
}
