package opencl64

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

// Set Bth to thermal noise (Brown).
// see temperature.cu
func SetTemperature(Bth, noise *data.Slice, k2mu0_Mu0VgammaDt float64, Msat, Temp, Alpha MSlice) {
	util.Argument(Bth.NComp() == 1 && noise.NComp() == 1)

	N := Bth.Len()
	cfg := make1DConf(N)

	Beff := (unsafe.Pointer)(nil)
	nois := (unsafe.Pointer)(nil)
	Msat_X := (unsafe.Pointer)(nil)
	Temp_X := (unsafe.Pointer)(nil)
	Alpha_X := (unsafe.Pointer)(nil)
	eventList := [](*cl.Event){}
	var tmpEvt *cl.Event

	if Bth != nil {
		Beff = Bth.DevPtr(0)
		tmpEvt = Bth.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	} else {
		panic("ERROR (SetTemperature): Bth pointer cannot be nil")
	}
	if noise != nil {
		nois = noise.DevPtr(0)
		tmpEvt = noise.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	} else {
		panic("ERROR (SetTemperature): Bth pointer cannot be nil")
	}
	if Msat.GetSlicePtr() != nil {
		Msat_X = Msat.DevPtr(0)
		tmpEvt = Msat.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if Temp.GetSlicePtr() != nil {
		Temp_X = Temp.DevPtr(0)
		tmpEvt = Temp.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if Alpha.GetSlicePtr() != nil {
		Alpha_X = Alpha.DevPtr(0)
		tmpEvt = Alpha.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_settemperature2_async(Beff, nois, float64(k2mu0_Mu0VgammaDt),
		Msat_X, Msat.Mul(0),
		Temp_X, Temp.Mul(0),
		Alpha_X, Alpha.Mul(0),
		N, cfg,
		eventList)

	Bth.SetEvent(0, event)

	glist := []GSlice{noise}
	if Msat_X != nil {
		glist = append(glist, Msat)
	}
	if Temp_X != nil {
		glist = append(glist, Temp)
	}
	if Alpha_X != nil {
		glist = append(glist, Alpha)
	}
	InsertEventIntoGSlices(event, glist)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in settemperature: %+v \n", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
		return
	}

	go WaitAndUpdateDataSliceEvents(event, glist, true)

}
