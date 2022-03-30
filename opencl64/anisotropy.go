package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

// Adds cubic anisotropy field to Beff.
func AddCubicAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, k3, c1, c2 MSlice) {
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	eventList := []*cl.Event{Beff.GetEvent(X), Beff.GetEvent(Y), Beff.GetEvent(Z),
		m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z)}

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
		N, cfg, eventList)

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addcubicanisotropy: %+v \n", err)
	}
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy2.cl
func AddUniaxialAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, u MSlice) {
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	eventList := []*cl.Event{Beff.GetEvent(X), Beff.GetEvent(Y), Beff.GetEvent(Z),
		m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z)}

	event := k_adduniaxialanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		k2.DevPtr(0), k2.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, eventList)

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addcubicanisotropy2: %+v \n", err)
	}
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy.cl
func AddUniaxialAnisotropy(Beff, m *data.Slice, Msat, k1, u MSlice) {
	util.Argument(Beff.Size() == m.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	eventList := []*cl.Event{Beff.GetEvent(X), Beff.GetEvent(Y), Beff.GetEvent(Z),
		m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z)}

	event := k_adduniaxialanisotropy_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, eventList)

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
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

	eventList := []*cl.Event{Beff.GetEvent(X), Beff.GetEvent(Y), Beff.GetEvent(Z),
		m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z)}

	event := k_addvoltagecontrolledanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		vcmaCoeff.DevPtr(0), vcmaCoeff.Mul(0),
		voltage.DevPtr(0), voltage.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg, eventList)

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addvoltagecontrolledanisotropy: %+v \n", err)
	}
}
