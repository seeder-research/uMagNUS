package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Landau-Lifshitz torque divided by gamma0:
//   - 1/(1+α²) [ m x B +  α m x (m x B) ]
//     torque in Tesla
//     m normalized
//     B in Tesla
// see lltorque.cl
func LLTorque(torque, m, B *data.Slice, alpha MSlice) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		lltorque__(torque, m, B, alpha, wg)
	} else {
		go lltorque__(torque, m, B, alpha, wg)
	}
	wg.Wait()
}

func lltorque__(torque, m, B *data.Slice, alpha MSlice, wg_ sync.WaitGroup) {
	torque.Lock(X)
	torque.Lock(Y)
	torque.Lock(Z)
	defer torque.Unlock(X)
	defer torque.Unlock(Y)
	defer torque.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)
	B.RLock(X)
	B.RLock(Y)
	B.RLock(Z)
	defer B.RUnlock(X)
	defer B.RUnlock(Y)
	defer B.RUnlock(Z)
	if alpha.GetSlicePtr() != nil {
		alpha.RLock()
		defer alpha.RUnlock()
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("lltorque failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := torque.Len()
	cfg := make1DConf(N)

	event := k_lltorque2_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		alpha.DevPtr(0), alpha.Mul(0), N, cfg, cmdqueue,
		nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in lltorque: %+v \n", err)
	}
}

// Landau-Lifshitz torque with precession disabled.
// Used by engine.Relax().
func LLNoPrecess(torque, m, B *data.Slice) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		llnoprocess__(torque, m, B, wg)
	} else {
		go llnoprocess__(torque, m, B, wg)
	}
	wg.Wait()
}

func llnoprocess__(torque, m, B *data.Slice, wg_ sync.WaitGroup) {
	torque.Lock(X)
	torque.Lock(Y)
	torque.Lock(Z)
	defer torque.Unlock(X)
	defer torque.Unlock(Y)
	defer torque.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)
	B.RLock(X)
	B.RLock(Y)
	B.RLock(Z)
	defer B.RUnlock(X)
	defer B.RUnlock(Y)
	defer B.RUnlock(Z)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("llnoprecess failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := torque.Len()
	cfg := make1DConf(N)

	event := k_llnoprecess_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z), N, cfg, cmdqueue,
		nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in llnoprecess: %+v \n", err)
	}
}
