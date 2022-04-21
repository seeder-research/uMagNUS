package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// shift dst by shx cells (positive or negative) along X-axis.
// new edge value is clampL at left edge or clampR at right edge.
func ShiftX(dst, src *data.Slice, shiftX int, clampL, clampR float32) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == src.Len())
	N := dst.Size()
	cfg := make3DConf(N)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent(0)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = src.GetEvent(0)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_shiftx_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftX, clampL, clampR, cfg,
		eventsList)

	dst.SetEvent(0, event)
	src.SetEvent(0, event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in shiftx failed: %+v \n", err)
		}
	}
}

func ShiftY(dst, src *data.Slice, shiftY int, clampL, clampR float32) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == src.Len())
	N := dst.Size()
	cfg := make3DConf(N)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent(0)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = src.GetEvent(0)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_shifty_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftY, clampL, clampR, cfg,
		eventsList)

	dst.SetEvent(0, event)
	src.SetEvent(0, event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in shifty failed: %+v \n", err)
		}
	}
}

func ShiftZ(dst, src *data.Slice, shiftZ int, clampL, clampR float32) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == src.Len())
	N := dst.Size()
	cfg := make3DConf(N)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent(0)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = src.GetEvent(0)
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_shiftz_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftZ, clampL, clampR, cfg,
		eventsList)

	dst.SetEvent(0, event)
	src.SetEvent(0, event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in shiftz failed: %+v \n", err)
		}
	}
}

// Like Shift, but for bytes
func ShiftBytes(dst, src *Bytes, m *data.Mesh, shiftX int, clamp byte) {
	N := m.Size()
	cfg := make3DConf(N)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent()
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = src.GetEvent()
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_shiftbytes_async(dst.Ptr, src.Ptr, N[X], N[Y], N[Z], shiftX, clamp, cfg, nil)

	dst.SetEvent(event)
	src.SetEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in shiftbytes failed: %+v \n", err)
		}
	}
}

func ShiftBytesY(dst, src *Bytes, m *data.Mesh, shiftY int, clamp byte) {
	N := m.Size()
	cfg := make3DConf(N)

	eventsList := []*cl.Event{}
	tmpEvt := dst.GetEvent()
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	tmpEvt = src.GetEvent()
	if tmpEvt != nil {
		eventsList = append(eventsList, tmpEvt)
	}
	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_shiftbytesy_async(dst.Ptr, src.Ptr, N[X], N[Y], N[Z], shiftY, clamp, cfg, nil)

	dst.SetEvent(event)
	src.SetEvent(event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in shiftbytesy failed: %+v \n", err)
		}
	}
}
