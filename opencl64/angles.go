package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

func SetPhi(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	event := k_setPhi_async(s.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y),
		N[X], N[Y], N[Z],
		cfg, []*cl.Event{s.GetEvent(0), m.GetEvent(X), m.GetEvent(Y)})

	s.SetEvent(X, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)

	// Force synchronization
	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in phi: %+v \n", err)
	}
	return
}

func SetTheta(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	event := k_setTheta_async(s.DevPtr(0), m.DevPtr(Z),
		N[X], N[Y], N[Z],
		cfg, []*cl.Event{s.GetEvent(0), m.GetEvent(Z)})

	s.SetEvent(X, event)
	m.SetEvent(Z, event)

	// Force synchronization
	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in theta: %+v \n", err)
	}
	return
}
