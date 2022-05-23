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
	tmpEvtL := torque.GetAllEvents(X)
	if len(tmpEvtL) > 0 {
		eventList = append(eventList, tmpEvtL...)
	}
	tmpEvtL = torque.GetAllEvents(Y)
	if len(tmpEvtL) > 0 {
		eventList = append(eventList, tmpEvtL...)
	}
	tmpEvtL = torque.GetAllEvents(Z)
	if len(tmpEvtL) > 0 {
		eventList = append(eventList, tmpEvtL...)
	}
	tmpEvt := m.GetEvent(X)
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

	glist := []GSlice{m}
	if J.GetSlicePtr() != nil {
		glist = append(glist, J)
	}
	if fixedP.GetSlicePtr != nil {
		glist = append(glist, fixedP)
	}
	if alpha.GetSlicePtr() != nil {
		glist = append(glist, alpha)
	}
	if ε_prime.GetSlicePtr() != nil {
		glist = append(glist, ε_prime)
	}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	if pol.GetSlicePtr() != nil {
		glist = append(glist, pol)
	}
	if λ.GetSlicePtr() != nil {
		glist = append(glist, λ)
	}
	if thickness.GetSlicePtr() != nil {
		glist = append(glist, thickness)
	}
	InsertEventIntoGSlices(event, glist)

	if Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in addslonczewskitorque2: %+v \n", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
		return
	}

	go WaitAndUpdateDataSliceEvents(event, glist, true)

}
