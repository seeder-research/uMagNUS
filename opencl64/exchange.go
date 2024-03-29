package opencl64

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
)

// Add exchange field to Beff.
// 	m: normalized magnetization
// 	B: effective field in Tesla
// 	Aex_red: Aex / (Msat * 1e18 m2)
// see exchange.cl
func AddExchange(B, m *data.Slice, Aex_red SymmLUT, Msat MSlice, regions *Bytes, mesh *data.Mesh) {
	c := mesh.CellSize()
	wx := float64(2 / (c[X] * c[X]))
	wy := float64(2 / (c[Y] * c[Y]))
	wz := float64(2 / (c[Z] * c[Z]))
	N := mesh.Size()
	pbc := mesh.PBC_code()
	cfg := make3DConf(N)

	eventList := []*cl.Event{}
	tmpEvtL := B.GetAllEvents(X)
	if len(tmpEvtL) > 0 {
		eventList = append(eventList, tmpEvtL...)
	}
	tmpEvtL = B.GetAllEvents(Y)
	if len(tmpEvtL) > 0 {
		eventList = append(eventList, tmpEvtL...)
	}
	tmpEvtL = B.GetAllEvents(Z)
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
	if Msat.GetSlicePtr() != nil {
		tmpEvt = Msat.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	tmpEvt = regions.GetEvent()
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_addexchange_async(B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		unsafe.Pointer(Aex_red), regions.Ptr,
		wx, wy, wz, N[X], N[Y], N[Z], pbc, cfg,
		eventList)

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
			fmt.Printf("WaitForEvents failed in addexchange: %+v", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
		regions.RemoveReadEvent(event)
		return
	}

	go WaitAndUpdateDataSliceEvents(event, glist, true)
	go func(ev *cl.Event, b *Bytes) {
		if err := cl.WaitForEvents([]*cl.Event{ev}); err != nil {
			fmt.Printf("WaitForEvents failed in addexchange: %+v \n", err)
		}
		b.RemoveReadEvent(ev)
	}(event, regions)

}

// Finds the average exchange strength around each cell, for debugging.
func ExchangeDecode(dst *data.Slice, Aex_red SymmLUT, regions *Bytes, mesh *data.Mesh) {
	c := mesh.CellSize()
	wx := float64(2 / (c[X] * c[X]))
	wy := float64(2 / (c[Y] * c[Y]))
	wz := float64(2 / (c[Z] * c[Z]))
	N := mesh.Size()
	pbc := mesh.PBC_code()
	cfg := make3DConf(N)

	eventList := []*cl.Event{}
	tmpEvtL := dst.GetAllEvents(0)
	if len(tmpEvtL) > 0 {
		eventList = append(eventList, tmpEvtL...)
	}
	tmpEvt := regions.GetEvent()
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_exchangedecode_async(dst.DevPtr(0), unsafe.Pointer(Aex_red), regions.Ptr,
		wx, wy, wz, N[X], N[Y], N[Z], pbc, cfg,
		eventList)

	dst.SetEvent(0, event)

	regions.InsertReadEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in exchangedecode: %+v", err)
		}
		regions.RemoveReadEvent(event)
		return
	}

	go func(ev *cl.Event, b *Bytes) {
		if err := cl.WaitForEvents([]*cl.Event{ev}); err != nil {
			fmt.Printf("WaitForEvents failed in exchangedecode: %+v \n", err)
		}
		b.RemoveReadEvent(ev)
	}(event, regions)

}
