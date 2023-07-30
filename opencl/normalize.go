package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Normalize vec to unit length, unless length or vol are zero.
func Normalize(vec, vol *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(vol == nil || vol.NComp() == 1)
	N := vec.Len()
	cfg := make1DConf(N)

	event := k_normalize2_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z), volPtr, N, cfg, ewl, q)

	vec.SetEvent(X, event)
	vec.SetEvent(Y, event)
	vec.SetEvent(Z, event)

	glist := []GSlice{}
	if vol != nil {
		glist = append(glist, vol)
		InsertEventIntoGSlices(event, glist)
	}

	if Synchronous || Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in normalize: %+v \n", err)
		}
	}

	return
}
