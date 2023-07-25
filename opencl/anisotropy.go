package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Adds cubic anisotropy field to Beff.
func AddCubicAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, k3, c1, c2 MSlice) {
	// synchronization should be done by code using
	// opencl library

	if Synchronous || Debug { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	// Checkout command queue from pool and launch kernel
	var addCubicAnisotropy2SyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &addCubicAnisotropy2SyncWaitGroup)
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
		N, cfg, nil,
		tmpQueue)

	// Check in command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &addCubicAnisotropy2SyncWaitGroup)
	ReturnQueuePool <- qwg

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)

	glist := []GSlice{m}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	if k1.GetSlicePtr() != nil {
		glist = append(glist, k1)
	}
	if k2.GetSlicePtr() != nil {
		glist = append(glist, k2)
	}
	if k3.GetSlicePtr() != nil {
		glist = append(glist, k3)
	}
	if c1.GetSlicePtr() != nil {
		glist = append(glist, c1)
	}
	if c2.GetSlicePtr() != nil {
		glist = append(glist, c2)
	}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug { // debug
		addCubicAnisotropy2SyncWaitGroup.Wait()
	}

	return
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy2.cl
func AddUniaxialAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, u MSlice) {
	// synchronization should be done by code using
	// opencl library

	if Synchronous || Debug { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	// Checkout command queue from pool and launch kernel
	var addUniaxialAnisotropy2SyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &addUniaxialAnisotropy2SyncWaitGroup)
	event := k_adduniaxialanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		k2.DevPtr(0), k2.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, nil,
		tmpQueue)

	// Check in command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &addUniaxialAnisotropy2SyncWaitGroup)
	ReturnQueuePool <- qwg

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)

	glist := []GSlice{m}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	if k1.GetSlicePtr() != nil {
		glist = append(glist, k1)
	}
	if k2.GetSlicePtr() != nil {
		glist = append(glist, k2)
	}
	if u.GetSlicePtr() != nil {
		glist = append(glist, u)
	}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug { // debug
		addUniaxialAnisotropy2SyncWaitGroup.Wait()
	}

	return
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy.cl
func AddUniaxialAnisotropy(Beff, m *data.Slice, Msat, k1, u MSlice) {
	// synchronization should be done by code using
	// opencl library

	if Synchronous || Debug { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	// Checkout command queue from pool and launch kernel
	var addUniaxialAnisotropySyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &addUniaxialAnisotropySyncWaitGroup)
	event := k_adduniaxialanisotropy_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, nil,
		tmpQueue)

	// Check in command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &addUniaxialAnisotropySyncWaitGroup)
	ReturnQueuePool <- qwg

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)

	glist := []GSlice{m}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	if k1.GetSlicePtr() != nil {
		glist = append(glist, k1)
	}
	if u.GetSlicePtr() != nil {
		glist = append(glist, u)
	}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug { // debug
		addUniaxialAnisotropySyncWaitGroup.Wait()
	}

	return
}

// Add voltage-conrtolled magnetic anisotropy field to Beff.
// see voltagecontrolledanisotropy2.cu
func AddVoltageControlledAnisotropy(Beff, m *data.Slice, Msat, vcmaCoeff, voltage, u MSlice) {
	// synchronization should be done by code using
	// opencl library

	if Synchronous || Debug { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	util.Argument(Beff.Size() == m.Size())

	checkSize(Beff, m, vcmaCoeff, voltage, u, Msat)

	N := Beff.Len()
	cfg := make1DConf(N)

	// Checkout command queue from pool and launch kernel
	var addVoltageControlledAnisotropySyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &addVoltageControlledAnisotropySyncWaitGroup)
	event := k_addvoltagecontrolledanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		vcmaCoeff.DevPtr(0), vcmaCoeff.Mul(0),
		voltage.DevPtr(0), voltage.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, nil,
		tmpQueue)

	// Check in command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &addVoltageControlledAnisotropy2SyncWaitGroup)
	ReturnQueuePool <- qwg

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)

	glist := []GSlice{m}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	if vcmaCoeff.GetSlicePtr() != nil {
		glist = append(glist, vcmaCoeff)
	}
	if voltage.GetSlicePtr() != nil {
		glist = append(glist, voltage)
	}
	if u.GetSlicePtr() != nil {
		glist = append(glist, u)
	}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug { // debug
		addCubicAnisotropy2SyncWaitGroup.Wait()
	}

	return
}
