package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Crop stores in dst a rectangle cropped from src at given offset position.
// dst size may be smaller than src.
func Crop(dst, src *data.Slice, offX, offY, offZ int) {
	D := dst.Size()
	S := src.Size()
	util.Argument(dst.NComp() == src.NComp())
	util.Argument(D[X]+offX <= S[X] && D[Y]+offY <= S[Y] && D[Z]+offZ <= S[Z])

	cfg := make3DConf(D)

	eventList := make([](*cl.Event), dst.NComp())
	for c := 0; c < dst.NComp(); c++ {
		eventWaitList := []*cl.Event{}
		tmpEvent := dst.GetEvent(c)
		if tmpEvent != nil {
			eventWaitList = append(eventWaitList, tmpEvent)
		}
		tmpEvent = src.GetEvent(c)
		if tmpEvent != nil {
			eventWaitList = append(eventWaitList, tmpEvent)
		}
		if len(eventWaitList) == 0 {
			eventWaitList = nil
		}
		eventList[c] = k_crop_async(dst.DevPtr(c), D[X], D[Y], D[Z],
			src.DevPtr(c), S[X], S[Y], S[Z],
			offX, offY, offZ, cfg, eventWaitList)
		dst.SetEvent(c, eventList[c])
		src.InsertReadEvent(c, eventList[c])
		go func(ev *cl.Event, idx int, sl *data.Slice) {
			if err := cl.WaitForEvents([]*cl.Event{ev}); err != nil {
				fmt.Printf("WaitForEvents failed in crop: %+v \n", err)
			}
			src.RemoveReadEvent(idx, ev)
		}(eventList[c], c, src)
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in crop: %+v \n", err)
		}
	}
}
