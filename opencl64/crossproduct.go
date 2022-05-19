package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

func CrossProduct(dst, a, b *data.Slice) {
	util.Argument(dst.NComp() == 3 && a.NComp() == 3 && b.NComp() == 3)
	util.Argument(dst.Len() == a.Len() && dst.Len() == b.Len())

	N := dst.Len()
	cfg := make1DConf(N)
	eventWaitList := []*cl.Event{}
	tmpEventL := dst.GetAllEvents(X)
	if len(tmpEventL) > 0 {
		eventWaitList = append(eventWaitList, tmpEventL...)
	}
	tmpEventL = dst.GetAllEvents(Y)
	if len(tmpEventL) > 0 {
		eventWaitList = append(eventWaitList, tmpEventL...)
	}
	tmpEventL = dst.GetAllEvents(Z)
	if len(tmpEventL) > 0 {
		eventWaitList = append(eventWaitList, tmpEventL...)
	}
	tmpEvent := a.GetEvent(X)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = a.GetEvent(Y)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = a.GetEvent(Z)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = b.GetEvent(X)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = b.GetEvent(Y)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = b.GetEvent(Z)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	if len(eventWaitList) == 0 {
		eventWaitList = nil
	}

	event := k_crossproduct_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		b.DevPtr(X), b.DevPtr(Y), b.DevPtr(Z),
		N, cfg, eventWaitList)

	dst.SetEvent(X, event)
	dst.SetEvent(Y, event)
	dst.SetEvent(Z, event)

	glist := []GSlice{a, b}
	InsertEventIntoGSlices(event, glist)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in crossproduct: %+v \n", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
		return
	}

	go WaitAndUpdateDataSliceEvents(event, glist, true)

}
