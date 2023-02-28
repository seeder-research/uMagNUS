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
		crossproduct__(dst, a, b, &wg)
	} else {
		go func() {
			crossproduct__(dst, a, b, &wg)
		}()
	}
	wg.Wait()
}

func crossproduct__(dst, a, b *data.Slice, wg_ *sync.WaitGroup) {
	dst.Lock(X)
	dst.Lock(Y)
	dst.Lock(Z)
	defer dst.Unlock(X)
	defer dst.Unlock(Y)
	defer dst.Unlock(Z)
	if dst.DevPtr(X) != a.DevPtr(X) {
		a.RLock(X)
		defer a.RUnlock(X)
	}
	if dst.DevPtr(Y) != a.DevPtr(Y) {
		a.RLock(Y)
		defer a.RUnlock(Y)
	}
	if dst.DevPtr(Z) != a.DevPtr(Z) {
		a.RLock(Z)
		defer a.RUnlock(Z)
	}
	if dst.DevPtr(X) != b.DevPtr(X) {
		b.RLock(X)
		defer b.RUnlock(X)
	}
	if dst.DevPtr(Y) != b.DevPtr(Y) {
		b.RLock(Y)
		defer b.RUnlock(Y)
	}
	if dst.DevPtr(Z) != b.DevPtr(Z) {
		b.RLock(Z)
		defer b.RUnlock(Z)
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("crossproduct failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_crossproduct_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		b.DevPtr(X), b.DevPtr(Y), b.DevPtr(Z),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in crossproduct: %+v \n", err)
	}
}
