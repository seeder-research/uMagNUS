package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// m = 1 / (4 + τ²(m x H)²) [{4 - τ²(m x H)²} m - 4τ(m x m x H)]
// note: torque from LLNoPrecess has negative sign
func Minimize(m, m0, torque *data.Slice, dt float32) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		minimize__(m, m0, torque, dt, wg)
	} else {
		go minimize__(m, m0, torque, dt, wg)
	}
	wg.Wait()
}

func minimize__(m, m0, torque *data.Slice, dt float32, wg_ sync.WaitGroup) {
	m.Lock(X)
	m.Lock(Y)
	m.Lock(Z)
	defer m.Unlock(X)
	defer m.Unlock(Y)
	defer m.Unlock(Z)
	m0.RLock(X)
	m0.RLock(Y)
	m0.RLock(Z)
	defer m0.RUnlock(X)
	defer m0.RUnlock(Y)
	defer m0.RUnlock(Z)
	torque.RLock(X)
	torque.RLock(Y)
	torque.RLock(Z)
	defer torque.RUnlock(X)
	defer torque.RUnlock(Y)
	defer torque.RUnlock(Z)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("minimize failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := m.Len()
	cfg := make1DConf(N)

	event := k_minimize_async(m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		m0.DevPtr(X), m0.DevPtr(Y), m0.DevPtr(Z),
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		dt, N, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in minimize: %+v \n", err)
	}
}
