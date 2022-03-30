package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Normalize vec to unit length, unless length or vol are zero.
func Normalize(vec, vol *data.Slice) {
	util.Argument(vol == nil || vol.NComp() == 1)
	N := vec.Len()
	cfg := make1DConf(N)

	eventList := []*cl.Event{vec.GetEvent(X), vec.GetEvent(Y), vec.GetEvent(Z)}
	volPtr := (unsafe.Pointer)(nil)

	if vol != nil {
		volPtr = vol.DevPtr(0)
		eventList = append(eventList, vol.GetEvent(0))
	}

	event := k_normalize2_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z), volPtr, N, cfg, eventList)

	vec.SetEvent(X, event)
	vec.SetEvent(Y, event)
	vec.SetEvent(Z, event)
	if vol != nil {
		vol.SetEvent(X, event)
	}
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in normalize: %+v \n", err)
	}
}
