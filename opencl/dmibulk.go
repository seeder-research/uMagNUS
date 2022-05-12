package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Add effective field due to bulk Dzyaloshinskii-Moriya interaction to Beff.
// See dmibulk.cl
func AddDMIBulk(Beff *data.Slice, m *data.Slice, Aex_red, D_red SymmLUT, Msat MSlice, regions *Bytes, mesh *data.Mesh, OpenBC bool) {
	cellsize := mesh.CellSize()
	N := Beff.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)
	var openBC byte
	if OpenBC {
		openBC = 1
	}

	eventWaitList := []*cl.Event{}
	tmpEvtL := Beff.GetAllEvents(X)
	if len(tmpEvtL) > 0 {
		eventWaitList = append(eventWaitList, tmpEvtL...)
	}
	tmpEvtL = Beff.GetAllEvents(Y)
	if len(tmpEvtL) > 0 {
		eventWaitList = append(eventWaitList, tmpEvtL...)
	}
	tmpEvtL = Beff.GetAllEvents(Z)
	if len(tmpEvtL) > 0 {
		eventWaitList = append(eventWaitList, tmpEvtL...)
	}
	tmpEvt := m.GetEvent(X)
	if tmpEvt != nil {
		eventWaitList = append(eventWaitList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Y)
	if tmpEvt != nil {
		eventWaitList = append(eventWaitList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Z)
	if tmpEvt != nil {
		eventWaitList = append(eventWaitList, tmpEvt)
	}
	if Msat.GetSlicePtr() != nil {
		tmpEvt = Msat.GetEvent(0)
		if tmpEvt != nil {
			eventWaitList = append(eventWaitList, tmpEvt)
		}
	}
	tmpEvt = regions.GetEvent()
	if tmpEvt != nil {
		eventWaitList = append(eventWaitList, tmpEvt)
	}
	if len(eventWaitList) == 0 {
		eventWaitList = nil
	}

	event := k_adddmibulk_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		unsafe.Pointer(Aex_red), unsafe.Pointer(D_red), regions.Ptr,
		float32(cellsize[X]), float32(cellsize[Y]), float32(cellsize[Z]),
		N[X], N[Y], N[Z], mesh.PBC_code(), openBC, cfg,
		eventWaitList)

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)

	glist := []GSlice{m}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	InsertEventIntoGSlices(event, glist)
	regions.InsertReadEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in adddmibulk: %+v \n", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
		regions.RemoveReadEvent(event)
		return
	}

	go WaitAndUpdateDataSliceEvents(event, glist, true)
	go func(ev *cl.Event, b *Bytes) {
		if err := cl.WaitForEvents([]*cl.Event{ev}); err != nil {
			fmt.Printf("WaitForEvents failed in adddmibulk: %+v \n", err)
		}
		b.RemoveReadEvent(ev)
	}(event, regions)

}
