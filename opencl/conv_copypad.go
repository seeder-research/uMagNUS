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
	// synchronization should be done by code using
	// opencl library

	if Synchronous || Debug { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}

	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Argument(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	cfg := make3DConf(dstsize)

	// Checkout command queue from pool and launch kernel
	var copyUnPadSyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &copyUnPadSyncWaitGroup)
	event := k_copyunpad_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], cfg, nil,
		tmpQueue)

	// Check in command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &copyUnPadSyncWaitGroup)
	ReturnQueuePool <- qwg

	dst.SetEvent(0, event)

	glist := []GSlice{src}
	InsertEventIntoGSlices(event, glist)

	if (Synchronous || Debug) {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in copyunpad: %+v \n", err)
		}
	}

	return
}

// Copies src into dst, which is larger, and multiplies by vol*Bsat.
// The remainder of dst is not filled with zeros.
// Used to zero-pad magnetization before convolution and in the meanwhile multiply m by its length.
func copyPadMul(dst, src, vol *data.Slice, dstsize, srcsize [3]int, Msat MSlice) {
	// synchronization should be done by code using
	// opencl library

	if Synchronous || Debug { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}

	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize))

	cfg := make3DConf(srcsize)

	// Checkout command queue from pool and launch kernel
	var copyPadMulSyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &copyPadMulSyncWaitGroup)
	event := k_copypadmul2_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z],
		Msat.DevPtr(0), Msat.Mul(0), vol.DevPtr(0), cfg, nil,
		tmpQueue)

	// Check in command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &copyPadMulSyncWaitGroup)
	ReturnQueuePool <- qwg

	dst.SetEvent(0, event)

	glist := []GSlice{src, vol}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	InsertEventIntoGSlices(event, glist)

	if (Synchronous || Debug) {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in copypadmul: %+v \n", err)
		}
	}

	return
}
