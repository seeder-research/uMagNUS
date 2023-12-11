package opencl

import (
	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

func SetPhi(s *data.Slice, m *data.Slice, queue *cl.CommandQueue, events []*cl.Event) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	k_setPhi_async(s.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y),
		N[X], N[Y], N[Z],
		cfg, queue, events)
}

func SetTheta(s *data.Slice, m *data.Slice, queue *cl.CommandQueue, events []*cl.Event) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	k_setTheta_async(s.DevPtr(0), m.DevPtr(Z),
		N[X], N[Y], N[Z],
		cfg, queue, events)
}
