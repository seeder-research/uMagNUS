package opencl

// Region paired spin torque calculations

import (
	"fmt"
	"math"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

func AddRegionSpinTorque(torque, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, J, alpha, pfix, pfree, λfix, λfree, ε_prime float32, mesh *data.Mesh) {
	c := mesh.CellSize()
	dX := float64(sX) * c[X]
	dY := float64(sY) * c[Y]
	dZ := float64(sZ) * c[Z]

	distsq := dX*dX + dY*dY + dZ*dZ
	cellwgt := math.Abs(dX*c[X]) + math.Abs(dY*c[Y]) + math.Abs(dZ*c[Z])
	if cellwgt > 0.0 {
		cellwgt = math.Sqrt(distsq) / cellwgt
	}

	N := mesh.Size()
	cfg := make3DConf(N)

	eventsList := []*cl.Event{}
	tmpEvtL := torque.GetAllEvents(X)
	if len(tmpEvtL) > 0 {
		eventsList = append(eventsList, tmpEvtL...)
	}
	tmpEvtL = torque.GetAllEvents(Y)
	if len(tmpEvtL) > 0 {
		eventsList = append(eventsList, tmpEvtL...)
	}
	tmpEvtL = torque.GetAllEvents(Z)
	if len(tmpEvtL) > 0 {
		eventsList = append(eventsList, tmpEvtL...)
	}
	tmpEvt := m.GetEvent(X)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Y)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Z)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if Msat.GetSlicePtr() != nil {
		tmpEvt = Msat.GetEvent(0)
		if tmpEvt != nil {
			eventsList = append(eventsList, tmpEvt)
		}
	}
	tmpEvt = regions.GetEvent()
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_addtworegionoommfslonczewskitorque_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, N[X], N[Y], N[Z],
		J, alpha, pfix, pfree, λfix, λfree, ε_prime, float32(cellwgt),
		cfg,
		eventsList)

	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)

	glist := []GSlice{m}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	InsertEventIntoGSlices(event, glist)
	regions.InsertReadEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in addtworegionoommfslonczewskitorque: %+v", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
		regions.RemoveReadEvent(event)
		return
	}

	go WaitAndUpdateDataSliceEvents(event, glist, true)
	go func(ev *cl.Event, b *Bytes) {
		if err := cl.WaitForEvents([]*cl.Event{ev}); err != nil {
			fmt.Printf("WaitForEvents failed in addtworegionoommfslonczewskitorque: %+v \n", err)
		}
		b.RemoveReadEvent(ev)
	}(event, regions)

}
