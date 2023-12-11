package opencl

import (
	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Normalize vec to unit length, unless length or vol are zero.
func Normalize(vec, vol *data.Slice, queue *cl.CommandQueue, events []*cl.Event) {
	util.Argument(vol == nil || vol.NComp() == 1)
	N := vec.Len()
	cfg := make1DConf(N)

	k_normalize2_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z),
		vol.DevPtr(0), N, cfg, queue, events)
}
