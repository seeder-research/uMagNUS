package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

func SetPhi(s *data.Slice, m *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	// Launch kernel
	event := k_setPhi_async(s.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y),
		N[X], N[Y], N[Z],
		cfg, ewl,
		q)

	if Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in setphi: %+v \n", err)
		}
	}

	return
}

func SetTheta(s *data.Slice, m *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	// Launch kernel
	event := k_setTheta_async(s.DevPtr(0), m.DevPtr(Z),
		N[X], N[Y], N[Z],
		cfg, ewl,
		q)

	if Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in settheta: %+v \n", err)
		}
	}

	return
}
