package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

// divide: dst[i] = a[i] / b[i]
// divide by zero automagically returns 0.0
func Divide(dst, a, b *data.Slice) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(a.Len() == N && a.NComp() == nComp && b.Len() == N && b.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		eventList[c] = k_divide_async(dst.DevPtr(c), a.DevPtr(c), b.DevPtr(c), N, cfg,
			[](*cl.Event){dst.GetEvent(c), a.GetEvent(c), b.GetEvent(c)})
	}
	err := cl.WaitForEvents(eventList)
	if err != nil {
		fmt.Printf("WaitForEvents failed in divide: %+v \n", err)
	}
}
