package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// dst += LUT[region], for vectors. Used to add terms to excitation.
func RegionAddV(dst *data.Slice, lut LUTPtrs, regions *Bytes) {
	util.Argument(dst.NComp() == 3)
	N := dst.Len()
	cfg := make1DConf(N)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent(X)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = dst.GetEvent(Y)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = dst.GetEvent(Z)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = regions.GetEvent()
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_regionaddv_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		lut[X], lut[Y], lut[Z], regions.Ptr, N, cfg, eventsList)

	dst.SetEvent(X, event)
	dst.SetEvent(Y, event)
	dst.SetEvent(Z, event)
	regions.SetEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in regionaddv failed: %+v \n", err)
		}
	}
}

// dst += LUT[region], for scalar. Used to add terms to scalar excitation.
func RegionAddS(dst *data.Slice, lut LUTPtr, regions *Bytes) {
	util.Argument(dst.NComp() == 1)
	N := dst.Len()
	cfg := make1DConf(N)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent(0)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = regions.GetEvent()
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_regionadds_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		eventsList)

	dst.SetEvent(0, event)
	regions.SetEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in regionadds failed: %+v \n", err)
		}
	}
}

// decode the regions+LUT pair into an uncompressed array
func RegionDecode(dst *data.Slice, lut LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent(0)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = regions.GetEvent()
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_regiondecode_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg,
		eventsList)

	dst.SetEvent(0, event)
	regions.SetEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in regiondecode failed: %+v \n", err)
		}
	}
}

// select the part of src within the specified region, set 0's everywhere else.
func RegionSelect(dst, src *data.Slice, regions *Bytes, region byte) {
	util.Argument(dst.NComp() == src.NComp())
	N := dst.Len()
	cfg := make1DConf(N)

	eventList := make([]*cl.Event, dst.NComp())
	for c := 0; c < dst.NComp(); c++ {
		intWaitList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intWaitList = append(intWaitList, tmpEvt)
		}
		tmpEvt = src.GetEvent(c)
		if tmpEvt != nil {
			intWaitList = append(intWaitList, tmpEvt)
		}
		if len(intWaitList) == 0 {
			intWaitList = nil
		}

		eventList[c] = k_regionselect_async(dst.DevPtr(c), src.DevPtr(c), regions.Ptr, region, N, cfg,
			intWaitList)

		dst.SetEvent(c, eventList[c])
		src.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents in regionselect failed: %+v \n", err)
		}
	}
}
