package opencl

// Region paired spin torque calculations

import (
	"fmt"
	"math"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

func AddRegionSpinTorque(torque, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, J, alpha, pfix, pfree, λfix, λfree, ε_prime float32, mesh *data.Mesh) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addregionspintorque__(torque, m, Msat, regions, regionA, regionB,
			sX, sY, sZ, J, alpha, pfix, pfree, λfix, λfree, ε_prime,
			mesh, wg)
	} else {
		go addregionspintorque__(torque, m, Msat, regions, regionA, regionB,
			sX, sY, sZ, J, alpha, pfix, pfree, λfix, λfree, ε_prime,
			mesh, wg)
	}
	wg.Wait()
}

func addregionspintorque__(torque, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, J, alpha, pfix, pfree, λfix, λfree, ε_prime float32, mesh *data.Mesh, wg_ sync.WaitGroup) {
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
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}
	if regions != nil {
		regions.RLock()
		defer regions.RUnlock()
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("addtworegionoommfslonczewskitorque failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	c := mesh.CellSize()
	dX := float64(sX) * c[X]
	dY := float64(sY) * c[Y]
	dZ := float64(sZ) * c[Z]

	distsq := dX*dX + dY*dY + dZ*dZ
	cellwgt := math.Abs(dX*c[X]) + math.Abs(dY*c[Y]) + math.Abs(dZ*c[Z])
	if cellwgt > 0.0 {
		cellwgt = math.Sqrt(distsq) / cellwgt
	}

	N := mesh.Size()
	cfg := make3DConf(N)

	event := k_addtworegionoommfslonczewskitorque_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, N[X], N[Y], N[Z],
		J, alpha, pfix, pfree, λfix, λfree, ε_prime, float32(cellwgt),
		cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addtworegionoommfslonczewskitorque: %+v", err)
	}
}
