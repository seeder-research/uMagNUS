package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Sets vector dst to zero where mask != 0.
func ZeroMask(dst *data.Slice, mask LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)

	eventList := make([]*cl.Event, dst.NComp())
	for c := 0; c < dst.NComp(); c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = regions.GetEvent()
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_zeromask_async(dst.DevPtr(c), unsafe.Pointer(mask), regions.Ptr, N, cfg, intEventList)

		dst.SetEvent(c, eventList[c])
		regions.SetEvent(eventList[c])
	}

	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in zeromask: %+v \n", err)
		}
	}
}
