package opencl

import (
	"fmt"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// dst += LUT[region], for vectors. Used to add terms to excitation.
func RegionAddV(dst *data.Slice, lut LUTPtrs, regions *Bytes) {
	util.Argument(dst.NComp() == 3)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		regionaddv__(dst, lut, regions, &wg)
	} else {
		go func() {
			regionaddv__(dst, lut, regions, &wg)
		}()
	}
	wg.Wait()
}

func regionaddv__(dst *data.Slice, lut LUTPtrs, regions *Bytes, wg_ *sync.WaitGroup) {
	dst.Lock(X)
	dst.Lock(Y)
	dst.Lock(Z)
	defer dst.Unlock(X)
	defer dst.Unlock(Y)
	defer dst.Unlock(Z)
	if regions != nil {
		regions.RLock()
		defer regions.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("regionaddv failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_regionaddv_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		lut[X], lut[Y], lut[Z], regions.Ptr, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in regionaddv failed: %+v \n", err)
	}
}

// dst += LUT[region], for scalar. Used to add terms to scalar excitation.
func RegionAddS(dst *data.Slice, lut LUTPtr, regions *Bytes) {
	util.Argument(dst.NComp() == 1)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		regionadds__(dst, lut, regions, &wg)
	} else {
		go func() {
			regionadds__(dst, lut, regions, &wg)
		}()
	}
	wg.Wait()
}

func regionadds__(dst *data.Slice, lut LUTPtr, regions *Bytes, wg_ *sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
	if regions != nil {
		regions.RLock()
		defer regions.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("regionadds failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_regionadds_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in regionadds failed: %+v \n", err)
	}
}

// decode the regions+LUT pair into an uncompressed array
func RegionDecode(dst *data.Slice, lut LUTPtr, regions *Bytes) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		regiondecode__(dst, lut, regions, &wg)
	} else {
		go func() {
			regiondecode__(dst, lut, regions, &wg)
		}()
	}
	wg.Wait()
}

func regiondecode__(dst *data.Slice, lut LUTPtr, regions *Bytes, wg_ *sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
	if regions != nil {
		regions.RLock()
		defer regions.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("addtworegionoommfslonczewskitorque failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_regiondecode_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in regiondecode failed: %+v \n", err)
	}
}

// select the part of src within the specified region, set 0's everywhere else.
func RegionSelect(dst, src *data.Slice, regions *Bytes, region byte) {
	util.Argument(dst.NComp() == src.NComp())

	var wg sync.WaitGroup
	numComp := dst.NComp()
	wg.Add(numComp)
	for c := 0; c < numComp; c++ {
		wg.Add(1)
		if Synchronous {
			regionselect__(dst, src, regions, region, c, &wg)
		} else {
			idx := c
			go func() {
				regionselect__(dst, src, regions, region, idx, &wg)
			}()
		}
	}
	wg.Wait()
}

func regionselect__(dst, src *data.Slice, regions *Bytes, region byte, c int, wg_ *sync.WaitGroup) {
	dst.Lock(c)
	defer dst.Unlock(c)
	src.RLock(c)
	defer src.RUnlock(c)
	if regions != nil {
		regions.RLock()
		defer regions.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("regionselect failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_regionselect_async(dst.DevPtr(c), src.DevPtr(c), regions.Ptr, region, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in regionselect: %+v \n", err)
	}
}
