package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
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

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		resize__(dst, src, layer, wg)
	} else {
		go resize__(dst, src, layer, wg)
	}
	wg.Wait()
}

func resize__(dst, src *data.Slice, layer int, wg_ sync.WaitGroup) {
	dstsize := dst.Size()
	scalex := srcsize[X] / dstsize[X]
	scaley := srcsize[Y] / dstsize[Y]

	dst.Lock(0)
	defer dst.Unlock(0)
	src.RLock(0)
	defer src.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("resize failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	cfg := make3DConf(dstsize)

	event := k_resize_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], layer, scalex, scaley, cfg,
		cmdqueue, nil)

	wg_.Done()

	// Synchronize for resize
	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in resize: %+v \n", err)
	}
}
