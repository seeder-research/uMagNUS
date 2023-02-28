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
		lltorque__(torque, m, B, alpha, &wg)
	} else {
		go func() {
			lltorque__(torque, m, B, alpha, &wg)
		}()
	}
	wg.Wait()
}

func lltorque__(torque, m, B *data.Slice, alpha MSlice, wg_ *sync.WaitGroup) {
	torque.Lock(X)
	torque.Lock(Y)
	torque.Lock(Z)
	defer torque.Unlock(X)
	defer torque.Unlock(Y)
	defer torque.Unlock(Z)
	if torque.DevPtr(X) != m.DevPtr(X) {
		m.RLock(X)
		defer m.RUnlock(X)
	}
	if torque.DevPtr(Y) != m.DevPtr(Y) {
		m.RLock(Y)
		defer m.RUnlock(Y)
	}
	if torque.DevPtr(Z) != m.DevPtr(Z) {
		m.RLock(Z)
		defer m.RUnlock(Z)
	}
	if torque.DevPtr(X) != B.DevPtr(X) {
		B.RLock(X)
		defer B.RUnlock(X)
	}
	if torque.DevPtr(Y) != B.DevPtr(Y) {
		B.RLock(Y)
		defer B.RUnlock(Y)
	}
	if torque.DevPtr(Z) != B.DevPtr(Z) {
		B.RLock(Z)
		defer B.RUnlock(Z)
	}
	if alpha.GetSlicePtr() != nil {
		alpha.RLock()
		defer alpha.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("lltorque failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := torque.Len()
	cfg := make1DConf(N)

	event := k_lltorque2_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		alpha.DevPtr(0), alpha.Mul(0), N, cfg, cmdqueue,
		nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in lltorque: %+v \n", err)
	}
}

// Landau-Lifshitz torque with precession disabled.
// Used by engine.Relax().
func LLNoPrecess(torque, m, B *data.Slice) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		llnoprocess__(torque, m, B, &wg)
	} else {
		go func() {
			llnoprocess__(torque, m, B, &wg)
		}()
	}
	wg.Wait()
}

func llnoprocess__(torque, m, B *data.Slice, wg_ *sync.WaitGroup) {
	torque.Lock(X)
	torque.Lock(Y)
	torque.Lock(Z)
	defer torque.Unlock(X)
	defer torque.Unlock(Y)
	defer torque.Unlock(Z)
	if torque.DevPtr(X) != m.DevPtr(X) {
		m.RLock(X)
		defer m.RUnlock(X)
	}
	if torque.DevPtr(Y) != m.DevPtr(Y) {
		m.RLock(Y)
		defer m.RUnlock(Y)
	}
	if torque.DevPtr(Z) != m.DevPtr(Z) {
		m.RLock(Z)
		defer m.RUnlock(Z)
	}
	if torque.DevPtr(X) != B.DevPtr(X) {
		B.RLock(X)
		defer B.RUnlock(X)
	}
	if torque.DevPtr(Y) != B.DevPtr(Y) {
		B.RLock(Y)
		defer B.RUnlock(Y)
	}
	if torque.DevPtr(Z) != B.DevPtr(Z) {
		B.RLock(Z)
		defer B.RUnlock(Z)
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("llnoprecess failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := torque.Len()
	cfg := make1DConf(N)

	event := k_llnoprecess_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z), N, cfg, cmdqueue,
		nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in llnoprecess: %+v \n", err)
	}
}
