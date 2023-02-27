package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Copies src (larger) into dst (smaller).
// Used to extract demag field after convolution on padded m.
func copyUnPad(dst, src *data.Slice, dstsize, srcsize [3]int) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Argument(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		copyunpad__(dst, src, dstsize, srcsize, &wg)
	} else {
		go copyunpad__(dst, src, dstsize, srcsize, &wg)
	}
	wg.Wait()
}

func copyunpad__(dst, src *data.Slice, dstsize, srcsize [3]int, wg_ *sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
	src.RLock(0)
	defer src.RUnlock(0)

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("copyunpad failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	cfg := make3DConf(dstsize)

	event := k_copyunpad_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], cfg,
		cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in copyunpad: %+v \n", err)
	}
}

// Copies src into dst, which is larger, and multiplies by vol*Bsat.
// The remainder of dst is not filled with zeros.
// Used to zero-pad magnetization before convolution and in the meanwhile multiply m by its length.
func copyPadMul(dst, src, vol *data.Slice, dstsize, srcsize [3]int, Msat MSlice) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		copypadmul__(dst, src, vol, dstsize, srcsize, Msat, &wg)
	} else {
		go copypadmul__(dst, src, vol, dstsize, srcsize, Msat, &wg)
	}
	wg.Wait()
}

func copypadmul__(dst, src, vol *data.Slice, dstsize, srcsize [3]int, Msat MSlice, wg_ *sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
	src.RLock(0)
	defer src.RUnlock(0)
	vol.RLock(0)
	defer vol.RUnlock(0)
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("copypadmul failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	cfg := make3DConf(srcsize)

	event := k_copypadmul2_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z],
		Msat.DevPtr(0), Msat.Mul(0), vol.DevPtr(0), cfg,
		cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in copypadmul: %+v \n", err)
	}
}
