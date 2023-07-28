package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

func SetPhi(s *data.Slice, m *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) {
	// need to synchronize on previous accesses to s and m
	// which can be seen from code using opencl library

	if Synchronous { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	// Launch kernel
	event := k_setPhi_async(s.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y),
		N[X], N[Y], N[Z],
		cfg, ewl,
		q)

	s.SetEvent(0, event)
	m.InsertReadEvent(X, event)
	m.InsertReadEvent(Y, event)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in setphi: %+v \n", err)
		}
	}

	return
}

func SetTheta(s *data.Slice, m *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) {
	// need to synchronize on previous accesses to s and m
	// which can be seen from code using opencl library

	if Synchronous { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	// Launch kernel
	event := k_setTheta_async(s.DevPtr(0), m.DevPtr(Z),
		N[X], N[Y], N[Z],
		cfg, ewl,
		q)

	s.SetEvent(0, event)
	m.InsertReadEvent(Z, event)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in settheta: %+v \n", err)
		}
	}

	return
}
