package opencl

import (
	"fmt"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Add effective field due to bulk Dzyaloshinskii-Moriya interaction to Beff.
// See dmibulk.cl
func AddDMIBulk(Beff *data.Slice, m *data.Slice, Aex_red, D_red SymmLUT, Msat MSlice, regions *Bytes, mesh *data.Mesh, OpenBC bool) {
	N := Beff.Size()
	util.Argument(m.Size() == N)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		adddmibulk__(Beff, m, Aex_red, D_red, Msat, regions, mesh, OpenBC, wg)
	} else {
		go adddmibulk__(Beff, m, Aex_red, D_red, Msat, regions, mesh, OpenBC, wg)
	}
	wg.Done()
}

func adddmibulk__(Beff *data.Slice, m *data.Slice, Aex_red, D_red SymmLUT, Msat MSlice, regions *Bytes, mesh *data.Mesh, OpenBC bool, wg_ sync.WaitGroup) {
	var openBC byte
	if OpenBC {
		openBC = 1
	}

	Beff.Lock(X)
	Beff.Lock(Y)
	Beff.Lock(Z)
	defer Beff.Unlock(X)
	defer Beff.Unlock(Y)
	defer Beff.Unlock(Z)
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
		fmt.Printf("adddmibulk failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	cellsize := mesh.CellSize()
	N := Beff.Size()
	cfg := make3DConf(N)

	event := k_adddmibulk_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		unsafe.Pointer(Aex_red), unsafe.Pointer(D_red), regions.Ptr,
		float32(cellsize[X]), float32(cellsize[Y]), float32(cellsize[Z]),
		N[X], N[Y], N[Z], mesh.PBC_code(), openBC, cfg, cmdqueue,
		nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in adddmibulk: %+v \n", err)
	}
}
