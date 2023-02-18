package opencl

import (
	"fmt"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Normalize vec to unit length, unless length or vol are zero.
func Normalize(vec, vol *data.Slice) {
	util.Argument(vol == nil || vol.NComp() == 1)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		normalize__(vec, vol, wg)
	} else {
		go normalize__(vec, vol, wg)
	}
	wg.Wait()
}

func normalize__(vec, vol *data.Slice, wg_ sync.WaitGroup) {
	N := vec.Len()
	cfg := make1DConf(N)

	vec.Lock(X)
	vec.Lock(Y)
	vec.Lock(Z)
	defer vec.Unlock(X)
	defer vec.Unlock(Y)
	defer vec.Unlock(Z)

	volPtr := (unsafe.Pointer)(nil)
	if vol != nil {
		volPtr = vol.DevPtr(0)
		vol.RLock(0)
		defer vol.RUnlock(0)
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("normalize failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	event := k_normalize2_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z),
		volPtr, N, cfg, cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in normalize: %+v \n", err)
	}
}
