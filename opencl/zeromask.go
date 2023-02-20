package opencl

import (
	"fmt"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Sets vector dst to zero where mask != 0.
func ZeroMask(dst *data.Slice, mask LUTPtr, regions *Bytes) {
	var wg sync.WaitGroup
	for c := 0; c < dst.NComp(); c++ {
		wg.Add(1)
		if Synchronous {
			zeromask__(dst, mask, regions, c, &wg)
		} else {
			go zeromask__(dst, mask, regions, c, &wg)
		}
	}
	wg.Wait()
}

func zeromask__(dst *data.Slice, mask LUTPtr, regions *Bytes, c int, wg_ *sync.WaitGroup) {
	dst.Lock(c)
	defer dst.Unlock(c)
	if regions != nil {
		regions.RLock()
		defer regions.RUnlock()
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("zeromask failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_zeromask_async(dst.DevPtr(c), unsafe.Pointer(mask), regions.Ptr, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in zeromask: %+v \n", err)
	}
}
