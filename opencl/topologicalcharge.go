package opencl

import (
	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Set s to the toplogogical charge density s = m · (∂m/∂x ❌ ∂m/∂y)
// see topologicalcharge.cl
func SetTopologicalCharge(s *data.Slice, m *data.Slice, mesh *data.Mesh, queue *cl.CommandQueue, events []*cl.Event) {
	cellsize := mesh.CellSize()
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)
	icxcy := float32(1.0 / (cellsize[X] * cellsize[Y]))

	k_settopologicalcharge_async(s.DevPtr(X),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		icxcy, N[X], N[Y], N[Z], mesh.PBC_code(), cfg,
		queue, events)
}
