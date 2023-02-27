package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Adds cubic anisotropy field to Beff.
func AddCubicAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, k3, c1, c2 MSlice) {
	util.Argument(Beff.Size() == m.Size())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addcubicanisotropy__(Beff, m, Msat, k1, k2, k3, c1, c2, &wg)
	} else {
		go addcubicanisotropy__(Beff, m, Msat, k1, k2, k3, c1, c2, &wg)
	}
	wg.Wait()
}

func addcubicanisotropy__(Beff, m *data.Slice, Msat, k1, k2, k3, c1, c2 MSlice, wg_ *sync.WaitGroup) {
	Beff.Lock(X)
	Beff.Lock(Y)
	Beff.Lock(Z)
	defer Beff.Unlock(X)
	defer Beff.Unlock(Y)
	defer Beff.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}
	if k1.GetSlicePtr() != nil {
		k1.RLock()
		defer k1.RUnlock()
	}
	if k2.GetSlicePtr() != nil {
		k2.RLock()
		defer k2.RUnlock()
	}
	if k3.GetSlicePtr() != nil {
		k3.RLock()
		defer k3.RUnlock()
	}
	if c1.GetSlicePtr() != nil {
		c1.RLock()
		defer c1.RUnlock()
	}
	if c2.GetSlicePtr() != nil {
		c2.RLock()
		defer c2.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("addcubicanisotropy2 failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := Beff.Len()
	cfg := make1DConf(N)

	event := k_addcubicanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		k2.DevPtr(0), k2.Mul(0),
		k3.DevPtr(0), k3.Mul(0),
		c1.DevPtr(X), c1.Mul(X),
		c1.DevPtr(Y), c1.Mul(Y),
		c1.DevPtr(Z), c1.Mul(Z),
		c2.DevPtr(X), c2.Mul(X),
		c2.DevPtr(Y), c2.Mul(Y),
		c2.DevPtr(Z), c2.Mul(Z),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in addcubicanisotropy: %+v \n", err)
	}
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy2.cl
func AddUniaxialAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, u MSlice) {
	util.Argument(Beff.Size() == m.Size())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		adduniaxialanisotropy2__(Beff, m, Msat, k1, k2, u, &wg)
	} else {
		go adduniaxialanisotropy2__(Beff, m, Msat, k1, k2, u, &wg)
	}
	wg.Wait()
}

func adduniaxialanisotropy2__(Beff, m *data.Slice, Msat, k1, k2, u MSlice, wg_ *sync.WaitGroup) {
	Beff.Lock(X)
	Beff.Lock(Y)
	Beff.Lock(Z)
	defer Beff.Unlock(X)
	defer Beff.Unlock(Y)
	defer Beff.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}
	if k1.GetSlicePtr() != nil {
		k1.RLock()
		defer k1.RUnlock()
	}
	if k2.GetSlicePtr() != nil {
		k2.RLock()
		defer k2.RUnlock()
	}
	if u.GetSlicePtr() != nil {
		u.RLock()
		defer u.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("adduniaxialanisotropy2 failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := Beff.Len()
	cfg := make1DConf(N)

	event := k_adduniaxialanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		k2.DevPtr(0), k2.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in adduniaxialanisotropy2: %+v \n", err)
	}
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy.cl
func AddUniaxialAnisotropy(Beff, m *data.Slice, Msat, k1, u MSlice) {
	util.Argument(Beff.Size() == m.Size())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		adduniaxialanisotropy__(Beff, m, Msat, k1, u, &wg)
	} else {
		go adduniaxialanisotropy__(Beff, m, Msat, k1, u, &wg)
	}
	wg.Wait()
}

func adduniaxialanisotropy__(Beff, m *data.Slice, Msat, k1, u MSlice, wg_ *sync.WaitGroup) {
	Beff.Lock(X)
	Beff.Lock(Y)
	Beff.Lock(Z)
	defer Beff.Unlock(X)
	defer Beff.Unlock(Y)
	defer Beff.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}
	if k1.GetSlicePtr() != nil {
		k1.RLock()
		defer k1.RUnlock()
	}
	if u.GetSlicePtr() != nil {
		u.RLock()
		defer u.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("adduniaxialanisotropy failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := Beff.Len()
	cfg := make1DConf(N)

	event := k_adduniaxialanisotropy_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("Wait for command to complete failed in adduniaxialanisotropy: %+v \n", err)
	}
}

// Add voltage-conrtolled magnetic anisotropy field to Beff.
// see voltagecontrolledanisotropy2.cu
func AddVoltageControlledAnisotropy(Beff, m *data.Slice, Msat, vcmaCoeff, voltage, u MSlice) {
	util.Argument(Beff.Size() == m.Size())

	checkSize(Beff, m, vcmaCoeff, voltage, u, Msat)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addvoltagecontrolledanisotropy__(Beff, m, Msat, vcmaCoeff, voltage, u, &wg)
	} else {
		go addvoltagecontrolledanisotropy__(Beff, m, Msat, vcmaCoeff, voltage, u, &wg)
	}
	wg.Wait()
}

func addvoltagecontrolledanisotropy__(Beff, m *data.Slice, Msat, vcmaCoeff, voltage, u MSlice, wg_ *sync.WaitGroup) {
	Beff.Lock(X)
	Beff.Lock(Y)
	Beff.Lock(Z)
	defer Beff.Unlock(X)
	defer Beff.Unlock(Y)
	defer Beff.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}
	if vcmaCoeff.GetSlicePtr() != nil {
		vcmaCoeff.RLock()
		defer vcmaCoeff.RUnlock()
	}
	if voltage.GetSlicePtr() != nil {
		voltage.RLock()
		defer voltage.RUnlock()
	}
	if u.GetSlicePtr() != nil {
		u.RLock()
		defer u.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("addvoltagecontrolledanisotropy failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := Beff.Len()
	cfg := make1DConf(N)

	event := k_addvoltagecontrolledanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		vcmaCoeff.DevPtr(0), vcmaCoeff.Mul(0),
		voltage.DevPtr(0), voltage.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("Wait for command to complete failed in addvoltagecontrolledanisotropy: %+v \n", err)
	}
}
