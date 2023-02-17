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
	cfg := make3DConf(N)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		setphi__(s, m, wg)
	} else {
		go setphi__(s, m, wg)
	}
	wg.Wait()
}

func setphi__(s *data.Slice, m *data.Slice, wg_ sync.WaitGroup) {
	s.Lock(X)
	defer s.Unlock(X)
	m.RLock(X)
	m.RLock(Y)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("phi failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	event := k_setPhi_async(s.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y),
		N[X], N[Y], N[Z],
		cfg, cmdqueue, nil)

	wg_.Done()

	// Force synchronization
	if err := cmdqueue.Finish(); err != nil {
		fmt.Printf("Wait for command to complete failed in phi: %+v \n", err)
	}
}

func SetTheta(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		settheta__(s, m, wg)
	} else {
		go settheta__(s, m, wg)
	}
	wg.Wait()
}

func settheta__(s *data.Slice, m *data.Slice, wg_ sync.WaitGroup) {
	s.Lock(X)
	defer s.Unlock(X)
	m.RLock(Z)
	defer m.RUnlock(Z)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("theta failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	event := k_setTheta_async(s.DevPtr(0), m.DevPtr(Z),
		N[X], N[Y], N[Z],
		cfg, cmdqueue, nil)

	wg_.Done()

	// Force synchronization
	if err := cmdqueue.Finish(); err != nil {
		fmt.Printf("Wait for command to complete failed in theta: %+v \n", err)
	}
}
