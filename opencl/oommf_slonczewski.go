package opencl

import (
	"fmt"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Add Slonczewski ST torque to torque (Tesla).
func AddOommfSlonczewskiTorque(torque, m *data.Slice, Msat, J, fixedP, alpha, pfix, pfree, λfix, λfree, ε_prime MSlice, mesh *data.Mesh) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addoommfslonczewskitorque__(torque, m, Msat, J,
			fixedP, alpha, pfix, pfree, λfix, λfree, ε_prime, mesh, wg)
	} else {
		go addoommfslonczewskitorque__(torque, m, Msat, J,
			fixedP, alpha, pfix, pfree, λfix, λfree, ε_prime, mesh, wg)
	}
	wg.Wait()
}

func addoommfslonczewskitorque__(torque, m *data.Slice, Msat, J, fixedP, alpha, pfix, pfree, λfix, λfree, ε_prime MSlice, mesh *data.Mesh, wg_ sync.WaitGroup) {
	N := torque.Len()
	cfg := make1DConf(N)
	flt := float32(mesh.WorldSize()[Z])

	torque.Lock(X)
	torque.Lock(Y)
	torque.Lock(Z)
	defer torque.Unlock(X)
	defer torque.Unlock(Y)
	defer torque.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)

	if J.GetSlicePtr() != nil {
		J.RLock()
		defer J.RUnlock()
	}
	if fixedP.GetSlicePtr() != nil {
		fixedP.RLock()
		defer fixedP.RUnlock()
	}
	if alpha.GetSlicePtr() != nil {
		alpha.RLock()
		defer alpha.RUnlock()
	}
	if ε_prime.GetSlicePtr() != nil {
		ε_prime.RLock()
		defer ε_prime.RUnlock()
	}
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}
	if pfix.GetSlicePtr() != nil {
		pfix.RLock()
		defer pfix.RUnlock()
	}
	if pfree.GetSlicePtr() != nil {
		pfree.RLock()
		defer pfree.RUnlock()
	}
	if λfix.GetSlicePtr() != nil {
		λfix.RLock()
		defer λfix.RUnlock()
	}
	if λfree.GetSlicePtr() != nil {
		λfree.RLock()
		defer λfree.RUnlock()
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("addoommfslonczewskitorque failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	event := k_addoommfslonczewskitorque_async(
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		J.DevPtr(Z), J.Mul(Z),
		fixedP.DevPtr(X), fixedP.Mul(X),
		fixedP.DevPtr(Y), fixedP.Mul(Y),
		fixedP.DevPtr(Z), fixedP.Mul(Z),
		alpha.DevPtr(0), alpha.Mul(0),
		pfix.DevPtr(0), pfix.Mul(0),
		pfree.DevPtr(0), pfree.Mul(0),
		λfix.DevPtr(0), λfix.Mul(0),
		λfree.DevPtr(0), λfree.Mul(0),
		ε_prime.DevPtr(0), ε_prime.Mul(0),
		unsafe.Pointer(uintptr(0)), flt,
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in addoommfslonczewskitorque: %+v \n", err)
	}
}
