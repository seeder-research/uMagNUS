package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Add Slonczewski ST torque to torque (Tesla).
func AddOommfSlonczewskiTorque(torque, m *data.Slice, Msat, J, fixedP, alpha, pfix, pfree, λfix, λfree, ε_prime MSlice, mesh *data.Mesh) {
	N := torque.Len()
	cfg := make1DConf(N)
	flt := float32(mesh.WorldSize()[Z])

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
	if pfix.GetSlicePtr() != nil {
		tmpEvt = pfix.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if pfree.GetSlicePtr() != nil {
		tmpEvt = pfree.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if λfix.GetSlicePtr() != nil {
		tmpEvt = λfix.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if λfree.GetSlicePtr() != nil {
		tmpEvt = λfree.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_addoommfslonczewskitorque_async(
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		J.DevPtr(Z), J.Mul(Z),
		fixedP.DevPtr(X), fixedP.Mul(X),
		fixedP.DevPtr(Y), fixedP.Mul(Y),
		fixedP.DevPtr(Z), fixedP.Mul(Z),
		alpha.DevPtr(0), alpha.Mul(0),
		pfix.DevPtr(0), pfix.Mul(0),
		pfree.DevPtr(0), pfree.Mul(0),
		λfix.DevPtr(0), λfix.Mul(0),
		λfree.DevPtr(0), λfree.Mul(0),
		ε_prime.DevPtr(0), ε_prime.Mul(0),
		unsafe.Pointer(uintptr(0)), flt,
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
	if Msat.GetSlicePtr() != nil {
		Msat.SetEvent(0, event)
	}
	if fixedP.GetSlicePtr() != nil {
		fixedP.SetEvent(X, event)
		fixedP.SetEvent(Y, event)
		fixedP.SetEvent(Z, event)
	}
	if ε_prime.GetSlicePtr() != nil {
		ε_prime.SetEvent(0, event)
	}
	if alpha.GetSlicePtr() != nil {
		alpha.SetEvent(0, event)
	}
	if pfix.GetSlicePtr() != nil {
		pfix.SetEvent(0, event)
	}
	if pfree.GetSlicePtr() != nil {
		pfree.SetEvent(0, event)
	}
	if λfix.GetSlicePtr() != nil {
		λfix.SetEvent(0, event)
	}
	if λfree.GetSlicePtr() != nil {
		λfree.SetEvent(0, event)
	}

	if Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in addoommfslonczewskitorque: %+v \n", err)
		}
	}
}
