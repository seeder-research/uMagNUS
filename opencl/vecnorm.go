package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// dst = sqrt(dot(a, a)),
func VecNorm(dst *data.Slice, a *data.Slice) {
	util.Argument(dst.NComp() == 1 && a.NComp() == 3)
	util.Argument(dst.Len() == a.Len())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		vecnorm__(dst, a, &wg)
	} else {
		go vecnorm__(dst, a, &wg)
	}
	wg.Wait()
}

func vecnorm__(dst *data.Slice, a *data.Slice, wg_ *sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
	a.RLock(X)
	a.RLock(Y)
	a.RLock(Z)
	defer a.RUnlock(X)
	defer a.RUnlock(Y)
	defer a.RUnlock(Z)

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("vecnorm failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_vecnorm_async(dst.DevPtr(0),
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in vecnorm: %+v \n", err)
	}
}
