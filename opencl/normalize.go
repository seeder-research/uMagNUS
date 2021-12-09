package opencl

import (
	"fmt"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/cl"
	"github.com/seeder-research/uMagNUS/util"
)

// Normalize vec to unit length, unless length or vol are zero.
func Normalize(vec, vol *data.Slice) {
	util.Argument(vol == nil || vol.NComp() == 1)
	N := vec.Len()
	cfg := make1DConf(N)
	var event *cl.Event

	if vol == nil {
		event = k_normalize2_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z), nil, N, cfg,
			[](*cl.Event){vec.GetEvent(X), vec.GetEvent(Y), vec.GetEvent(Z)})
	} else {
		event = k_normalize2_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z), vol.DevPtr(0), N, cfg,
			[](*cl.Event){vec.GetEvent(X), vec.GetEvent(Y), vec.GetEvent(Z), vol.GetEvent(X)})
	}

	vec.SetEvent(X, event)
	vec.SetEvent(Y, event)
	vec.SetEvent(Z, event)
	if vol != nil {
		vol.SetEvent(X, event)
	}
	err := cl.WaitForEvents([]*cl.Event{event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in normalize: %+v \n", err)
	}
}
