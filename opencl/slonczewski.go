package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Add Slonczewski ST torque to torque (Tesla).
// see slonczewski.cl
func AddSlonczewskiTorque2(torque, m *data.Slice, Msat, J, fixedP, alpha, pol, λ, ε_prime, thickness MSlice, flp float64, mesh *data.Mesh) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addslonczewskitorque2__(torque, m, Msat, J,
			fixedP, alpha, pol, λ, ε_prime, thickness, flp, mesh, &wg)
	} else {
		go addslonczewskitorque2__(torque, m, Msat, J,
			fixedP, alpha, pol, λ, ε_prime, thickness, flp, mesh, &wg)
	}
	wg.Wait()
}

func addslonczewskitorque2__(torque, m *data.Slice, Msat, J, fixedP, alpha, pol, λ, ε_prime, thickness MSlice, flp float64, mesh *data.Mesh, wg_ *sync.WaitGroup) {
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
	if pol.GetSlicePtr() != nil {
		pol.RLock()
		defer pol.RUnlock()
	}
	if λ.GetSlicePtr() != nil {
		λ.RLock()
		defer λ.RUnlock()
	}
	if thickness.GetSlicePtr() != nil {
		thickness.RLock()
		defer thickness.RUnlock()
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("addslonczewskitorque2 failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := torque.Len()
	cfg := make1DConf(N)
	meshThickness := mesh.WorldSize()[Z]

	event := k_addslonczewskitorque2_async(
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		J.DevPtr(Z), J.Mul(Z),
		fixedP.DevPtr(X), fixedP.Mul(X),
		fixedP.DevPtr(Y), fixedP.Mul(Y),
		fixedP.DevPtr(Z), fixedP.Mul(Z),
		alpha.DevPtr(0), alpha.Mul(0),
		pol.DevPtr(0), pol.Mul(0),
		λ.DevPtr(0), λ.Mul(0),
		ε_prime.DevPtr(0), ε_prime.Mul(0),
		thickness.DevPtr(0), thickness.Mul(0),
		float32(meshThickness),
		float32(flp),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in addslonczewskitorque2: %+v \n", err)
	}
}
