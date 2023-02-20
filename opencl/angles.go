package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

func SetPhi(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		setphi__(s, m, &wg)
	} else {
		go setphi__(s, m, &wg)
	}
	wg.Wait()
}

func setphi__(s *data.Slice, m *data.Slice, wg_ *sync.WaitGroup) {
	s.Lock(X)
	defer s.Unlock(X)
	m.RLock(X)
	m.RLock(Y)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("setphi failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := s.Size()
	cfg := make3DConf(N)

	event := k_setPhi_async(s.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y),
		N[X], N[Y], N[Z],
		cfg, cmdqueue, nil)

	wg_.Done()

	// Force synchronization
	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in setphi: %+v \n", err)
	}
}

func SetTheta(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		settheta__(s, m, &wg)
	} else {
		go settheta__(s, m, &wg)
	}
	wg.Wait()
}

func settheta__(s *data.Slice, m *data.Slice, wg_ *sync.WaitGroup) {
	s.Lock(X)
	defer s.Unlock(X)
	m.RLock(Z)
	defer m.RUnlock(Z)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("settheta failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := s.Size()
	cfg := make3DConf(N)

	event := k_setTheta_async(s.DevPtr(0), m.DevPtr(Z),
		N[X], N[Y], N[Z],
		cfg, cmdqueue, nil)

	wg_.Done()

	// Force synchronization
	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in settheta: %+v \n", err)
	}
}
