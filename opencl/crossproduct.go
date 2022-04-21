package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

func CrossProduct(dst, a, b *data.Slice) {
	util.Argument(dst.NComp() == 3 && a.NComp() == 3 && b.NComp() == 3)
	util.Argument(dst.Len() == a.Len() && dst.Len() == b.Len())

	N := dst.Len()
	cfg := make1DConf(N)
	eventWaitList := []*cl.Event{}
	tmpEvent := dst.GetEvent(X)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = dst.GetEvent(Y)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = dst.GetEvent(Z)
	if tmpEvent != nil {
		eventWaitList = append(eventWaitList, tmpEvent)
	}
	tmpEvent = a.GetEvent(X)
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
	a.SetEvent(X, event)
	a.SetEvent(Y, event)
	a.SetEvent(Z, event)
	b.SetEvent(X, event)
	b.SetEvent(Y, event)
	b.SetEvent(Z, event)
	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in crossproduct: %+v \n", err)
		}
	}
}
