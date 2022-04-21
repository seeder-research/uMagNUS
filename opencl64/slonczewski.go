package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
)

// Add Slonczewski ST torque to torque (Tesla).
// see slonczewski.cl
func AddSlonczewskiTorque2(torque, m *data.Slice, Msat, J, fixedP, alpha, pol, λ, ε_prime MSlice, thickness MSlice, flp float64, mesh *data.Mesh) {
	N := torque.Len()
	cfg := make1DConf(N)
	meshThickness := mesh.WorldSize()[Z]

	eventList := [](*cl.Event){}
	tmpEvt := torque.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = torque.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = torque.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	if J.GetSlicePtr() != nil {
		tmpEvt = J.GetEvent(Z)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if fixedP.GetSlicePtr() != nil {
		tmpEvt = fixedP.GetEvent(X)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
		tmpEvt = fixedP.GetEvent(Y)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
		tmpEvt = fixedP.GetEvent(Z)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if alpha.GetSlicePtr() != nil {
		tmpEvt = alpha.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if ε_prime.GetSlicePtr() != nil {
		tmpEvt = ε_prime.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if Msat.GetSlicePtr() != nil {
		tmpEvt = Msat.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if pol.GetSlicePtr() != nil {
		tmpEvt = pol.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if λ.GetSlicePtr() != nil {
		tmpEvt = λ.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if thickness.GetSlicePtr() != nil {
		tmpEvt = thickness.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_addslonczewskitorque2_async(
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		J.DevPtr(Z), J.Mul(Z),
		fixedP.DevPtr(X), fixedP.Mul(X),
		fixedP.DevPtr(Y), fixedP.Mul(Y),
		fixedP.DevPtr(Z), fixedP.Mul(Z),
		alpha.DevPtr(0), alpha.Mul(0),
		pol.DevPtr(0), pol.Mul(0),
		λ.DevPtr(0), λ.Mul(0),
		ε_prime.DevPtr(0), ε_prime.Mul(0),
		thickness.DevPtr(0), thickness.Mul(0),
		float64(meshThickness),
		float64(flp),
		N, cfg, eventList)

	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	if J.GetSlicePtr() != nil {
		J.SetEvent(Z, event)
	}
	if fixedP.GetSlicePtr != nil {
		fixedP.SetEvent(X, event)
		fixedP.SetEvent(Y, event)
		fixedP.SetEvent(Z, event)
	}
	if alpha.GetSlicePtr() != nil {
		alpha.SetEvent(0, event)
	}
	if ε_prime.GetSlicePtr() != nil {
		ε_prime.SetEvent(0, event)
	}
	if Msat.GetSlicePtr() != nil {
		Msat.SetEvent(0, event)
	}
	if pol.GetSlicePtr() != nil {
		pol.SetEvent(0, event)
	}
	if λ.GetSlicePtr() != nil {
		λ.SetEvent(0, event)
	}
	if thickness.GetSlicePtr() != nil {
		thickness.SetEvent(0, event)
	}

	if Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in addslonczewskitorque2: %+v \n", err)
		}
	}
}
