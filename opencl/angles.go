package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

func SetPhi(s *data.Slice, m *data.Slice) {
	// need to synchronize on previous accesses to s and m
	// which can be seen from code using opencl library

	if Synchronous { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	// Checkout command queue from pool and launch kernel
	var setPhiSyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &setPhiSyncWaitGroup)
	event := k_setPhi_async(s.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y),
		N[X], N[Y], N[Z],
		cfg, nil,
		tmpQueue)

	// Check in command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &setPhiSyncWaitGroup)
	ReturnQueuePool <- qwg

	s.SetEvent(0, event)
	m.InsertReadEvent(X, event)
	m.InsertReadEvent(Y, event)

	return
}

func SetTheta(s *data.Slice, m *data.Slice) {
	// need to synchronize on previous accesses to s and m
	// which can be seen from code using opencl library

	if Synchronous { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)

	// Checkout command queue from pool and launch kernel
	var setThetaSyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &setThetaSyncWaitGroup)
	event := k_setTheta_async(s.DevPtr(0), m.DevPtr(Z),
		N[X], N[Y], N[Z],
		cfg, nil,
		tmpQueue)

	// Check in command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &setPhiSyncWaitGroup)
	ReturnQueuePool <- qwg

	s.SetEvent(0, event)
	m.InsertReadEvent(Z, event)

	return
}
