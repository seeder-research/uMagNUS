package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// divide: dst[i] = a[i] / b[i]
// divide by zero automagically returns 0.0
func Divide(dst, a, b *data.Slice) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(a.Len() == N && a.NComp() == nComp && b.Len() == N && b.NComp() == nComp)

	var wg sync.WaitGroup
	for c := 0; c < nComp; c++ {
		wg.Add(1)
		if Synchronous {
			divide__(dst, a, b, c, wg)
		} else {
			go divide__(dst, a, b, c, wg)
		}
	}
	wg.Wait()
}

func divide__(dst, a, b *data.Slice, idx int, wg_ sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	a.RLock(idx)
	defer a.RUnlock(idx)
	b.RLock(idx)
	defer b.RUnlock(idx)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("divide failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := dst.Len()
	nComp := dst.NComp()
	cfg := make1DConf(N)

	ev := k_divide_async(dst.DevPtr(idx), a.DevPtr(idx), b.DevPtr(idx), N, cfg, cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{ev}); err != nil {
		fmt.Printf("WaitForEvents failed in divide: %+v \n", err)
	}
}
