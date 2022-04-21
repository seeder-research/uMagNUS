package opencl64

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
)

// SetMaxAngle sets dst to the maximum angle of each cells magnetization with all of its neighbors,
// provided the exchange stiffness with that neighbor is nonzero.
func SetMaxAngle(dst, m *data.Slice, Aex_red SymmLUT, regions *Bytes, mesh *data.Mesh) {
	N := mesh.Size()
	pbc := mesh.PBC_code()
	cfg := make3DConf(N)

	eventList := [](*cl.Event){}
	tmpEvt := dst.GetEvent(0)
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
	tmpEvt = regions.GetEvent()
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_setmaxangle_async(dst.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		unsafe.Pointer(Aex_red), regions.Ptr,
		N[X], N[Y], N[Z], pbc, cfg, eventList)

	dst.SetEvent(0, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	regions.SetEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in setmaxangle: %+v \n", err)
		}
	}
}
