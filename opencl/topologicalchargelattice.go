package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Topological charge according to Berg and Lüscher
func SetTopologicalChargeLattice(s *data.Slice, m *data.Slice, mesh *data.Mesh) {
	N := s.Size()
	util.Argument(m.Size() == N)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		settopologicalcharglattice__(s, m, mesh, wg)
	} else {
		go settopologicalcharglattice__(s, m, mesh, wg)
	}
	wg.Wait()
}

func settopologicalcharglattice__(s *data.Slice, m *data.Slice, mesh *data.Mesh, wg_ sync.WaitGroup) {
	s.Lock(0)
	defer s.Unlock(0)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("settopologicalchargelattice failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	cellsize := mesh.CellSize()
	N := s.Size()
	cfg := make3DConf(N)
	icxcy := float32(1.0 / (cellsize[X] * cellsize[Y]))

	event := k_settopologicalchargelattice_async(
		s.DevPtr(X),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		icxcy, N[X], N[Y], N[Z], mesh.PBC_code(),
		cfg, cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in settopologicalchargelattice: %+v \n", err)
	}
}
