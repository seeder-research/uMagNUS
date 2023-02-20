package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// dst += prefactor * dot(a, b), as used for energy density
func AddDotProduct(dst *data.Slice, prefactor float32, a, b *data.Slice) {
	util.Argument(dst.NComp() == 1 && a.NComp() == 3 && b.NComp() == 3)
	util.Argument(dst.Len() == a.Len() && dst.Len() == b.Len())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		dotproduct__(dst, prefactor, a, b, &wg)
	} else {
		go dotproduct__(dst, prefactor, a, b, &wg)
	}
	wg.Wait()
}

func dotproduct__(dst *data.Slice, prefactor float32, a, b *data.Slice, wg_ *sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
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
		fmt.Printf("adddotproduct failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_dotproduct_async(dst.DevPtr(0), prefactor,
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		b.DevPtr(X), b.DevPtr(Y), b.DevPtr(Z),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in adddotproduct: %+v \n", err)
	}
}
