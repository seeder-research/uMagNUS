package opencl

import (
	"fmt"
	"unsafe"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl/cl"
	"github.com/seeder-research/uMagNUS/util"
)

// Adds cubic anisotropy field to Beff.
func AddCubicAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, k3, c1, c2 MSlice) {
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	Beff_X := (unsafe.Pointer)(nil)
	Beff_Y := (unsafe.Pointer)(nil)
	Beff_Z := (unsafe.Pointer)(nil)
	m_X := (unsafe.Pointer)(nil)
	m_Y := (unsafe.Pointer)(nil)
	m_Z := (unsafe.Pointer)(nil)
	Msat_X := (unsafe.Pointer)(nil)
	k1Ptr := (unsafe.Pointer)(nil)
	k2Ptr := (unsafe.Pointer)(nil)
	k3Ptr := (unsafe.Pointer)(nil)
	c1X := (unsafe.Pointer)(nil)
	c1Y := (unsafe.Pointer)(nil)
	c1Z := (unsafe.Pointer)(nil)
	c2X := (unsafe.Pointer)(nil)
	c2Y := (unsafe.Pointer)(nil)
	c2Z := (unsafe.Pointer)(nil)
	eventList := [](*cl.Event){}

	if Beff != nil {
		Beff_X = Beff.DevPtr(X)
		Beff_Y = Beff.DevPtr(Y)
		Beff_Z = Beff.DevPtr(Z)
		eventList = append(eventList, Beff.GetEvent(X))
		eventList = append(eventList, Beff.GetEvent(Y))
		eventList = append(eventList, Beff.GetEvent(Z))
	}
	if m != nil {
		m_X = m.DevPtr(X)
		m_Y = m.DevPtr(Y)
		m_Z = m.DevPtr(Z)
		eventList = append(eventList, m.GetEvent(X))
		eventList = append(eventList, m.GetEvent(Y))
		eventList = append(eventList, m.GetEvent(Z))
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat_X = Msat.DevPtr(0)
		eventList = append(eventList, Msat.GetEvent(0))
	}
	if k1.GetSlicePtr(0) != nil {
		k1Ptr = k1.DevPtr(0)
		eventList = append(eventList, k1.GetEvent(0))
	}
	if k2.GetSlicePtr(0) != nil {
		k2Ptr = k2.DevPtr(0)
		eventList = append(eventList, k2.GetEvent(0))
	}
	if k3.GetSlicePtr(0) != nil {
		k3Ptr = k3.DevPtr(0)
		eventList = append(eventList, k3.GetEvent(0))
	}
	if c1.GetSlicePtr(X) != nil {
		c1X = c1.DevPtr(X)
		eventList = append(eventList, c1.GetEvent(X))
	}
	if c1.GetSlicePtr(Y) != nil {
		c1Y = c1.DevPtr(Y)
		eventList = append(eventList, c1.GetEvent(Y))
	}
	if c1.GetSlicePtr(Z) != nil {
		c1Z = c1.DevPtr(Z)
		eventList = append(eventList, c1.GetEvent(Z))
	}
	if c2.GetSlicePtr(X) != nil {
		c2X = c2.DevPtr(X)
		eventList = append(eventList, c2.GetEvent(X))
	}
	if c2.GetSlicePtr(Y) != nil {
		c2Y = c2.DevPtr(Y)
		eventList = append(eventList, c2.GetEvent(Y))
	}
	if c2.GetSlicePtr(Z) != nil {
		c2Z = c2.DevPtr(Z)
		eventList = append(eventList, c2.GetEvent(Z))
	}

	event := k_addcubicanisotropy2_async(
		Beff_X, Beff_Y, Beff_Z,
		m_X, m_Y, m_Z,
		Msat_X, Msat.Mul(0),
		k1Ptr, k1.Mul(0),
		k2Ptr, k2.Mul(0),
		k3Ptr, k3.Mul(0),
		c1X, c1.Mul(X),
		c1Y, c1.Mul(Y),
		c1Z, c1.Mul(Z),
		c2X, c2.Mul(X),
		c2Y, c2.Mul(Y),
		c2Z, c2.Mul(Z),
		N, cfg, eventList)

	if Beff != nil {
		Beff.SetEvent(X, event)
		Beff.SetEvent(Y, event)
		Beff.SetEvent(Z, event)
	}
	if m != nil {
		m.SetEvent(X, event)
		m.SetEvent(Y, event)
		m.SetEvent(Z, event)
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat.SetEvent(0, event)
	}
	if k1.GetSlicePtr(0) != nil {
		k1.SetEvent(0, event)
	}
	if k2.GetSlicePtr(0) != nil {
		k2.SetEvent(0, event)
	}
	if k3.GetSlicePtr(0) != nil {
		k3.SetEvent(0, event)
	}
	if c1.GetSlicePtr(X) != nil {
		c1.SetEvent(X, event)
	}
	if c1.GetSlicePtr(Y) != nil {
		c1.SetEvent(Y, event)
	}
	if c1.GetSlicePtr(Z) != nil {
		c1.SetEvent(Z, event)
	}
	if c2.GetSlicePtr(X) != nil {
		c2.SetEvent(X, event)
	}
	if c2.GetSlicePtr(Y) != nil {
		c2.SetEvent(Y, event)
	}
	if c2.GetSlicePtr(Z) != nil {
		c2.SetEvent(Z, event)
	}
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in addcubicanisotropy: %+v \n", err)
	}
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy2.cl
func AddUniaxialAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, u MSlice) {
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	Beff_X := (unsafe.Pointer)(nil)
	Beff_Y := (unsafe.Pointer)(nil)
	Beff_Z := (unsafe.Pointer)(nil)
	m_X := (unsafe.Pointer)(nil)
	m_Y := (unsafe.Pointer)(nil)
	m_Z := (unsafe.Pointer)(nil)
	Msat_X := (unsafe.Pointer)(nil)
	k1Ptr := (unsafe.Pointer)(nil)
	k2Ptr := (unsafe.Pointer)(nil)
	uX := (unsafe.Pointer)(nil)
	uY := (unsafe.Pointer)(nil)
	uZ := (unsafe.Pointer)(nil)
	eventList := [](*cl.Event){}
	if Beff != nil {
		Beff_X = Beff.DevPtr(X)
		Beff_Y = Beff.DevPtr(Y)
		Beff_Z = Beff.DevPtr(Z)
		eventList = append(eventList, Beff.GetEvent(X))
		eventList = append(eventList, Beff.GetEvent(Y))
		eventList = append(eventList, Beff.GetEvent(Z))
	}
	if m != nil {
		m_X = m.DevPtr(X)
		m_Y = m.DevPtr(Y)
		m_Z = m.DevPtr(Z)
		eventList = append(eventList, m.GetEvent(X))
		eventList = append(eventList, m.GetEvent(Y))
		eventList = append(eventList, m.GetEvent(Z))
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat_X = Msat.DevPtr(0)
		eventList = append(eventList, Msat.GetEvent(0))
	}
	if k1.GetSlicePtr(0) != nil {
		k1Ptr = k1.DevPtr(0)
		eventList = append(eventList, k1.GetEvent(0))
	}
	if k2.GetSlicePtr(0) != nil {
		k2Ptr = k2.DevPtr(0)
		eventList = append(eventList, k2.GetEvent(0))
	}
	if u.GetSlicePtr(X) != nil {
		uX = u.DevPtr(X)
		eventList = append(eventList, u.GetEvent(X))
	}
	if u.GetSlicePtr(Y) != nil {
		uY = u.DevPtr(Y)
		eventList = append(eventList, u.GetEvent(Y))
	}
	if u.GetSlicePtr(Z) != nil {
		uZ = u.DevPtr(Z)
		eventList = append(eventList, u.GetEvent(Z))
	}

	event := k_adduniaxialanisotropy2_async(
		Beff_X, Beff_Y, Beff_Z,
		m_X, m_Y, m_Z,
		Msat_X, Msat.Mul(0),
		k1Ptr, k1.Mul(0),
		k2Ptr, k2.Mul(0),
		uX, u.Mul(X),
		uY, u.Mul(Y),
		uZ, u.Mul(Z),
		N, cfg, eventList)

	if Beff != nil {
		Beff.SetEvent(X, event)
		Beff.SetEvent(Y, event)
		Beff.SetEvent(Z, event)
	}
	if m != nil {
		m.SetEvent(X, event)
		m.SetEvent(Y, event)
		m.SetEvent(Z, event)
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat.SetEvent(0, event)
	}
	if k1.GetSlicePtr(0) != nil {
		k1.SetEvent(0, event)
	}
	if k2.GetSlicePtr(0) != nil {
		k2.SetEvent(0, event)
	}
	if u.GetSlicePtr(X) != nil {
		u.SetEvent(X, event)
	}
	if u.GetSlicePtr(Y) != nil {
		u.SetEvent(Y, event)
	}
	if u.GetSlicePtr(Z) != nil {
		u.SetEvent(Z, event)
	}
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in addcubicanisotropy2: %+v \n", err)
	}
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy.cl
func AddUniaxialAnisotropy(Beff, m *data.Slice, Msat, k1, u MSlice) {
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	Beff_X := (unsafe.Pointer)(nil)
	Beff_Y := (unsafe.Pointer)(nil)
	Beff_Z := (unsafe.Pointer)(nil)
	m_X := (unsafe.Pointer)(nil)
	m_Y := (unsafe.Pointer)(nil)
	m_Z := (unsafe.Pointer)(nil)
	Msat_X := (unsafe.Pointer)(nil)
	k1Ptr := (unsafe.Pointer)(nil)
	uX := (unsafe.Pointer)(nil)
	uY := (unsafe.Pointer)(nil)
	uZ := (unsafe.Pointer)(nil)
	eventList := [](*cl.Event){}
	if Beff != nil {
		Beff_X = Beff.DevPtr(X)
		Beff_Y = Beff.DevPtr(Y)
		Beff_Z = Beff.DevPtr(Z)
		eventList = append(eventList, Beff.GetEvent(X))
		eventList = append(eventList, Beff.GetEvent(Y))
		eventList = append(eventList, Beff.GetEvent(Z))
	}
	if m != nil {
		m_X = m.DevPtr(X)
		m_Y = m.DevPtr(Y)
		m_Z = m.DevPtr(Z)
		eventList = append(eventList, m.GetEvent(X))
		eventList = append(eventList, m.GetEvent(Y))
		eventList = append(eventList, m.GetEvent(Z))
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat_X = Msat.DevPtr(0)
		eventList = append(eventList, Msat.GetEvent(0))
	}
	if k1.GetSlicePtr(0) != nil {
		k1Ptr = k1.DevPtr(0)
		eventList = append(eventList, k1.GetEvent(0))
	}
	if u.GetSlicePtr(X) != nil {
		uX = u.DevPtr(X)
		eventList = append(eventList, u.GetEvent(X))
	}
	if u.GetSlicePtr(Y) != nil {
		uY = u.DevPtr(Y)
		eventList = append(eventList, u.GetEvent(Y))
	}
	if u.GetSlicePtr(Z) != nil {
		uZ = u.DevPtr(Z)
		eventList = append(eventList, u.GetEvent(Z))
	}

	event := k_adduniaxialanisotropy_async(
		Beff_X, Beff_Y, Beff_Z,
		m_X, m_Y, m_Z,
		Msat_X, Msat.Mul(0),
		k1Ptr, k1.Mul(0),
		uX, u.Mul(X),
		uY, u.Mul(Y),
		uZ, u.Mul(Z),
		N, cfg, eventList)

	if Beff != nil {
		Beff.SetEvent(X, event)
		Beff.SetEvent(Y, event)
		Beff.SetEvent(Z, event)
	}
	if m != nil {
		m.SetEvent(X, event)
		m.SetEvent(Y, event)
		m.SetEvent(Z, event)
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat.SetEvent(0, event)
	}
	if k1.GetSlicePtr(0) != nil {
		k1.SetEvent(0, event)
	}
	if u.GetSlicePtr(X) != nil {
		u.SetEvent(X, event)
	}
	if u.GetSlicePtr(Y) != nil {
		u.SetEvent(Y, event)
	}
	if u.GetSlicePtr(Z) != nil {
		u.SetEvent(Z, event)
	}
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in addcubicanisotropy: %+v \n", err)
	}
}

// Add voltage-conrtolled magnetic anisotropy field to Beff.
// see voltagecontrolledanisotropy2.cu
func AddVoltageControlledAnisotropy(Beff, m *data.Slice, Msat, vcmaCoeff, voltage, u MSlice) {
	util.Argument(Beff.Size() == m.Size())

	checkSize(Beff, m, vcmaCoeff, voltage, u, Msat)

	N := Beff.Len()
	cfg := make1DConf(N)

	Beff_X := (unsafe.Pointer)(nil)
	Beff_Y := (unsafe.Pointer)(nil)
	Beff_Z := (unsafe.Pointer)(nil)
	m_X := (unsafe.Pointer)(nil)
	m_Y := (unsafe.Pointer)(nil)
	m_Z := (unsafe.Pointer)(nil)
	Msat_X := (unsafe.Pointer)(nil)
	vcmaPtr := (unsafe.Pointer)(nil)
	voltPtr := (unsafe.Pointer)(nil)
	uX := (unsafe.Pointer)(nil)
	uY := (unsafe.Pointer)(nil)
	uZ := (unsafe.Pointer)(nil)
	eventList := [](*cl.Event){}
	if Beff != nil {
		Beff_X = Beff.DevPtr(X)
		Beff_Y = Beff.DevPtr(Y)
		Beff_Z = Beff.DevPtr(Z)
		eventList = append(eventList, Beff.GetEvent(X))
		eventList = append(eventList, Beff.GetEvent(Y))
		eventList = append(eventList, Beff.GetEvent(Z))
	}
	if m != nil {
		m_X = m.DevPtr(X)
		m_Y = m.DevPtr(Y)
		m_Z = m.DevPtr(Z)
		eventList = append(eventList, m.GetEvent(X))
		eventList = append(eventList, m.GetEvent(Y))
		eventList = append(eventList, m.GetEvent(Z))
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat_X = Msat.DevPtr(0)
		eventList = append(eventList, Msat.GetEvent(0))
	}
	if vcmaCoeff.GetSlicePtr(0) != nil {
		vcmaPtr = vcmaCoeff.DevPtr(0)
		eventList = append(eventList, vcmaCoeff.GetEvent(0))
	}
	if voltage.GetSlicePtr(0) != nil {
		voltPtr = voltage.DevPtr(0)
		eventList = append(eventList, voltage.GetEvent(0))
	}
	if u.GetSlicePtr(X) != nil {
		uX = u.DevPtr(X)
		eventList = append(eventList, u.GetEvent(X))
	}
	if u.GetSlicePtr(Y) != nil {
		uY = u.DevPtr(Y)
		eventList = append(eventList, u.GetEvent(Y))
	}
	if u.GetSlicePtr(Z) != nil {
		uZ = u.DevPtr(Z)
		eventList = append(eventList, u.GetEvent(Z))
	}

	event := k_addvoltagecontrolledanisotropy2_async(
		Beff_X, Beff_Y, Beff_Z,
		m_X, m_Y, m_Z,
		Msat_X, Msat.Mul(0),
		vcmaPtr, vcmaCoeff.Mul(0),
		voltPtr, voltage.Mul(0),
		uX, u.Mul(X),
		uY, u.Mul(Y),
		uZ, u.Mul(Z),
		N, cfg, eventList)

	if Beff != nil {
		Beff.SetEvent(X, event)
		Beff.SetEvent(Y, event)
		Beff.SetEvent(Z, event)
	}
	if m != nil {
		m.SetEvent(X, event)
		m.SetEvent(Y, event)
		m.SetEvent(Z, event)
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat.SetEvent(0, event)
	}
	if vcmaCoeff.GetSlicePtr(0) != nil {
		vcmaCoeff.SetEvent(0, event)
	}
	if voltage.GetSlicePtr(0) != nil {
		voltage.SetEvent(0, event)
	}
	if u.GetSlicePtr(X) != nil {
		u.SetEvent(X, event)
	}
	if u.GetSlicePtr(Y) != nil {
		u.SetEvent(Y, event)
	}
	if u.GetSlicePtr(Z) != nil {
		u.SetEvent(Z, event)
	}
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in addvoltagecontrolledanisotropy: %+v \n", err)
	}
}
