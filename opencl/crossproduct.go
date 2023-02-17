package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

func CrossProduct(dst, a, b *data.Slice) {
	util.Argument(dst.NComp() == 3 && a.NComp() == 3 && b.NComp() == 3)
	util.Argument(dst.Len() == a.Len() && dst.Len() == b.Len())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		crossproduct__(dst, a, b, wg)
	} else {
		go crossproduct__(dst, a, b, wg)
	}
	wg.Wait()
}

func crossproduct__(dst, a, b *data.Slice, wg_ sync.WaitGroup) {
	dst.Lock(X)
	dst.Lock(Y)
	dst.Lock(Z)
	defer dst.Unlock(X)
	defer dst.Unlock(Y)
	defer dst.Unlock(Z)
	a.RLock(X)
	a.RLock(Y)
	a.RLock(Z)
	defer a.RUnlock(X)
	defer a.RUnlock(Y)
	defer a.RUnlock(Z)
	b.RLock(X)
	b.RLock(Y)
	b.RLock(Z)
	defer b.RUnlock(X)
	defer b.RUnlock(Y)
	defer b.RUnlock(Z)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("crossproduct failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_crossproduct_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		b.DevPtr(X), b.DevPtr(Y), b.DevPtr(Z),
		N, cfg, cmdqueue, nil)

	wg.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in crossproduct: %+v \n", err)
	}
}
