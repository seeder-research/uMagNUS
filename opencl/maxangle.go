package opencl

import (
	"fmt"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// SetMaxAngle sets dst to the maximum angle of each cells magnetization with all of its neighbors,
// provided the exchange stiffness with that neighbor is nonzero.
func SetMaxAngle(dst, m *data.Slice, Aex_red SymmLUT, regions *Bytes, mesh *data.Mesh) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		maxangle__(dst, m, Aex_red, regions, mesh, wg)
	} else {
		go maxangle__(dst, m, Aex_red, regions, mesh, wg)
	}
	wg.Wait()
}

func maxangle__(dst, m *data.Slice, Aex_red SymmLUT, regions *Bytes, mesh *data.Mesh, wg_ sync.WaitGroup) {
	dst.Lock(X)
	defer dst.Unlock(X)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)
	regions.RLock()
	defer regions.RUnlock()

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("setmaxangle failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := mesh.Size()
	pbc := mesh.PBC_code()
	cfg := make3DConf(N)

	event := k_setmaxangle_async(dst.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		unsafe.Pointer(Aex_red), regions.Ptr,
		N[X], N[Y], N[Z], pbc, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in setmaxangle: %+v \n", err)
	}
}
