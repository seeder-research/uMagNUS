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

	if Bth != nil {
		Beff = Bth.DevPtr(0)
		eventList = append(eventList, Bth.GetEvent(0))
	}
	if noise != nil {
		nois = noise.DevPtr(0)
		eventList = append(eventList, noise.GetEvent(0))
	}
	if Msat.GetSlicePtr(0) != nil {
		Msat_X = Msat.DevPtr(0)
		eventList = append(eventList, Msat.GetEvent(0))
	}
	if Temp.GetSlicePtr(0) != nil {
		Temp_X = Temp.DevPtr(0)
		eventList = append(eventList, Temp.GetEvent(0))
	}
	if Alpha.GetSlicePtr(0) != nil {
		Alpha_X = Alpha.DevPtr(0)
		eventList = append(eventList, Alpha.GetEvent(0))
	}
	event := k_settemperature2_async(Beff, nois, float64(k2mu0_Mu0VgammaDt),
		Msat_X, Msat.Mul(0),
		Temp_X, Temp.Mul(0),
		Alpha_X, Alpha.Mul(0),
		N, cfg,
		eventList)

	if Beff != nil {
		Bth.SetEvent(0, event)
	}
	if nois != nil {
		noise.SetEvent(0, event)
	}
	if Msat_X != nil {
		Msat.SetEvent(0, event)
	}
	if Temp_X != nil {
		Temp.SetEvent(0, event)
	}
	if Alpha_X != nil {
		Alpha.SetEvent(0, event)
	}
	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in settemperature: %+v \n", err)
	}
}
