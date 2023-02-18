package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Add Zhang-Li ST torque (Tesla) to torque.
// see zhangli2.cl
func AddZhangLiTorque(torque, m *data.Slice, Msat, J, alpha, xi, pol MSlice, mesh *data.Mesh) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addzhanglitorque__(torque, m, Msat, J, alpha, xi, pol, mesh, wg)
	} else {
		go addzhanglitorque__(torque, m, Msat, J, alpha, xi, pol, mesh, wg)
	}
	wg.Wait()
}

func addzhanglitorque__(torque, m *data.Slice, Msat, J, alpha, xi, pol MSlice, mesh *data.Mesh, wg_ sync.WaitGroup) {
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
	if alpha.GetSlicePtr() != nil {
		alpha.RLock()
		defer alpha.RUnlock()
	}
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}
	if pol.GetSlicePtr() != nil {
		pol.RLock()
		defer pol.RUnlock()
	}
	if xi.GetSlicePtr() != nil {
		xi.RLock()
		defer xi.RUnlock()
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("addzhanglitorque failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	c := mesh.CellSize()
	N := mesh.Size()
	cfg := make3DConf(N)

	event := k_addzhanglitorque2_async(
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		J.DevPtr(X), J.Mul(X),
		J.DevPtr(Y), J.Mul(Y),
		J.DevPtr(Z), J.Mul(Z),
		alpha.DevPtr(0), alpha.Mul(0),
		xi.DevPtr(0), xi.Mul(0),
		pol.DevPtr(0), pol.Mul(0),
		float32(c[X]), float32(c[Y]), float32(c[Z]),
		N[X], N[Y], N[Z], mesh.PBC_code(), cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addzhanglitorque: %+v \n", err)
	}
}
