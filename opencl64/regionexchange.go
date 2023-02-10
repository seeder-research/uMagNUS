package opencl64

// Region paired exchange interaction

import (
	"fmt"
	"math"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
)

// Add exchange field to Beff.
//	m: normalized magnetization
//	B: effective field in Tesla
func AddRegionExchangeField(B, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, sig, sig2 float64, mesh *data.Mesh) {
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
	tmpEvtL := B.GetAllEvents(X)
	if len(tmpEvtL) > 0 {
		eventsList = append(eventsList, tmpEvtL...)
	}
	tmpEvtL = B.GetAllEvents(Y)
	if len(tmpEvtL) > 0 {
		eventsList = append(eventsList, tmpEvtL...)
	}
	tmpEvtL = B.GetAllEvents(Z)
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
	if len(eventsList) == 0 {
		eventsList = nil
	}

	sig_eff := sig * float64(cellwgt)
	sig2_eff := sig2 * float64(cellwgt)

	event := k_tworegionexchange_field_async(B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, sig_eff, sig2_eff, N[X], N[Y], N[Z], cfg,
		eventsList)

	B.SetEvent(X, event)
	B.SetEvent(Y, event)
	B.SetEvent(Z, event)

	glist := []GSlice{m}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	InsertEventIntoGSlices(event, glist)
	regions.InsertReadEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in addtworegionexchange_field: %+v", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
		regions.RemoveReadEvent(event)
		return
	}

	go WaitAndUpdateDataSliceEvents(event, glist, true)
	go func(ev *cl.Event, b *Bytes) {
		if err := cl.WaitForEvents([]*cl.Event{ev}); err != nil {
			fmt.Printf("WaitForEvents failed in addtworegionexchange_field: %+v \n", err)
		}
		b.RemoveReadEvent(ev)
	}(event, regions)

}

func AddRegionExchangeEdens(Edens, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, sig, sig2 float64, mesh *data.Mesh) {
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
	tmpEvtL := Edens.GetAllEvents(0)
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

	sig_eff := sig * float64(cellwgt)
	sig2_eff := sig2 * float64(cellwgt)

	event := k_tworegionexchange_edens_async(Edens.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, sig_eff, sig2_eff, N[X], N[Y], N[Z], cfg,
		eventsList)

	Edens.SetEvent(0, event)

	glist := []GSlice{m}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	InsertEventIntoGSlices(event, glist)
	regions.InsertReadEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in addtworegionexchange_edens: %+v", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
		regions.RemoveReadEvent(event)
		return
	}

	go WaitAndUpdateDataSliceEvents(event, glist, true)
	go func(ev *cl.Event, b *Bytes) {
		if err := cl.WaitForEvents([]*cl.Event{ev}); err != nil {
			fmt.Printf("WaitForEvents failed in addtworegionexchange_edens: %+v \n", err)
		}
		b.RemoveReadEvent(ev)
	}(event, regions)

}
