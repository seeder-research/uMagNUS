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

	eventList := []*cl.Event{}
	tmpEvt := vec.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = vec.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = vec.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	volPtr := (unsafe.Pointer)(nil)
	if vol != nil {
		volPtr = vol.DevPtr(0)
		eventList = append(eventList, vol.GetEvent(0))
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_normalize2_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z), volPtr, N, cfg, eventList)

	vec.SetEvent(X, event)
	vec.SetEvent(Y, event)
	vec.SetEvent(Z, event)
	if vol != nil {
		vol.SetEvent(0, event)
	}

	if Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in normalize: %+v \n", err)
		}
	}
}
