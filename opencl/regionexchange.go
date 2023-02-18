package opencl

// Region paired exchange interaction

import (
	"fmt"
	"math"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Add exchange field to Beff.
//	m: normalized magnetization
//	B: effective field in Tesla
func AddRegionExchangeField(B, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, sig, sig2 float32, mesh *data.Mesh) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addregionexchangefield__(B, m, Msat, regions, regionA, regionB,
			sX, sY, sZ, sig, sig2, mesh, wg)
	} else {
		go addregionexchangefield__(B, m, Msat, regions, regionA, regionB,
			sX, sY, sZ, sig, sig2, mesh, wg)
	}
	wg.Wait()
}

func addregionexchangefield__(B, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, sig, sig2 float32, mesh *data.Mesh, wg_ sync.WaitGroup) {
	B.Lock(X)
	B.Lock(Y)
	B.Lock(Z)
	defer B.Unlock(X)
	defer B.Unlock(Y)
	defer B.Unlock(Z)
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
		regions.Rlock()
		defer regions.RUnlock()
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("addtworegionexchange_field failed to create command queue: %+v \n", err)
		return nil
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

	sig_eff := sig * float32(cellwgt)
	sig2_eff := sig2 * float32(cellwgt)

	event := k_tworegionexchange_field_async(B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, sig_eff, sig2_eff, N[X], N[Y], N[Z], cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addtworegionexchange_field: %+v", err)
	}
}

func AddRegionExchangeEdens(Edens, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, sig, sig2 float32, mesh *data.Mesh) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addregionexchangeedens__(Edens, m, Msat, regions, regionA, regionB,
			sX, sY, sZ, sig, sig2, mesh, wg)
	} else {
		go addregionexchangeedens__(Edens, m, Msat, regions, regionA, regionB,
			sX, sY, sZ, sig, sig2, mesh, wg)
	}
	wg.Wait()
}

func addregionexchangeedens__(Edens, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, sig, sig2 float32, mesh *data.Mesh, wg_ sync.WaitGroup) {
	Edens.Lock(0)
	defer Edens.Unlock(0)
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
		fmt.Printf("addtworegionexchange_edens failed to create command queue: %+v \n", err)
		return nil
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

	sig_eff := sig * float32(cellwgt)
	sig2_eff := sig2 * float32(cellwgt)

	event := k_tworegionexchange_edens_async(Edens.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, sig_eff, sig2_eff, N[X], N[Y], N[Z], cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addtworegionexchange_edens: %+v", err)
	}
}
