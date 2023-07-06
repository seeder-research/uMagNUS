package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

func SetPhi(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	waitEvents := []*cl.Event{m.GetEvent(X), m.GetEvent(Y)}
	tmpEvtL := s.GetAllEvents(0)
	if len(tmpEvtL) > 0 {
		waitEvents = append(waitEvents, tmpEvtL...)
	}
	event := k_setPhi_async(s.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y),
		N[X], N[Y], N[Z],
		cfg, waitEvents)

	s.SetEvent(0, event)
	m.InsertReadEvent(X, event)
	m.InsertReadEvent(Y, event)

	// Force synchronization TODO: needed??
	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in phi: %+v \n", err)
	}
	m.RemoveReadEvent(X, event)
	m.RemoveReadEvent(Y, event)
	return
}

func SetTheta(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	waitEvents := []*cl.Event{m.GetEvent(Z)}
	tmpEvtL := s.GetAllEvents(0)
	if len(tmpEvtL) > 0 {
		waitEvents = append(waitEvents, tmpEvtL...)
	}
	event := k_setTheta_async(s.DevPtr(0), m.DevPtr(Z),
		N[X], N[Y], N[Z],
		cfg, waitEvents)

	s.SetEvent(0, event)
	m.InsertReadEvent(Z, event)

	// Force synchronization TODO: needed??
	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in theta: %+v \n", err)
	}
	m.RemoveReadEvent(Z, event)
	return
}
