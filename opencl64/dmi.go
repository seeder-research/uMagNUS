package opencl64

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

// Add effective field of Dzyaloshinskii-Moriya interaction to Beff (Tesla).
// According to Bagdanov and Röβler, PRL 87, 3, 2001. eq.8 (out-of-plane symmetry breaking).
// See dmi.cl
func AddDMI(Beff *data.Slice, m *data.Slice, Aex_red, Dex_red SymmLUT, Msat MSlice, regions *Bytes, mesh *data.Mesh, OpenBC bool) {
	cellsize := mesh.CellSize()
	N := Beff.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)
	var openBC byte
	if OpenBC {
		openBC = 1
	}

	eventWaitList := []*cl.Event{}
	tmpEvent := Beff.GetEvent(X)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = Beff.GetEvent(Y)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = Beff.GetEvent(Z)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = m.GetEvent(X)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = m.GetEvent(Y)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = m.GetEvent(Z)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = regions.GetEvent()
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	if len(eventWaitList) == 0 {
		eventWaitList = nil
	}
	event := k_adddmi_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		unsafe.Pointer(Aex_red), unsafe.Pointer(Dex_red), regions.Ptr,
		float64(cellsize[X]), float64(cellsize[Y]), float64(cellsize[Z]),
		N[X], N[Y], N[Z], mesh.PBC_code(), openBC, cfg,
		eventWaitList)

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	regions.SetEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in adddmi: %+v \n", err)
		}
	}
}
