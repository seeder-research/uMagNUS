package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
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
		waitList := []*cl.Event{}
		tmpEvtL := dst.GetAllEvents(c)
		if len(tmpEvtL) > 0 {
			waitList = append(waitList, tmpEvtL...)
		}
		tmpEvt := a.GetEvent(c)
		if tmpEvt != nil {
			waitList = append(waitList, tmpEvt)
		}
		tmpEvt = b.GetEvent(c)
		if tmpEvt != nil {
			waitList = append(waitList, tmpEvt)
		}
		if len(waitList) == 0 {
			waitList = nil
		}
		eventList[c] = k_divide_async(dst.DevPtr(c), a.DevPtr(c), b.DevPtr(c), N, cfg,
			waitList)
		dst.SetEvent(c, eventList[c])
		a.InsertReadEvent(c, eventList[c])
		b.InsertReadEvent(c, eventList[c])
		go func(ev *cl.Event, idx int, sl []*data.Slice) {
			if err := cl.WaitForEvents([]*cl.Event{ev}); err != nil {
				fmt.Printf("WaitForEvents failed in crop: %+v \n", err)
			}
			for _, sa := range sl {
				sa.RemoveReadEvent(idx, ev)
			}
		}(eventList[c], c, []*data.Slice{a, b})
	}
	if err := cl.WaitForEvents(eventList); err != nil {
		fmt.Printf("WaitForEvents failed in divide: %+v \n", err)
	}
}
