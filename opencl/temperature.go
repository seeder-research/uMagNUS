package opencl

import (
	"fmt"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Set Bth to thermal noise (Brown).
// see temperature.cu
func SetTemperature(Bth, noise *data.Slice, k2mu0_Mu0VgammaDt float64, Msat, Temp, Alpha MSlice) {
	util.Argument(Bth.NComp() == 1 && noise.NComp() == 1)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		settemperature__(Bth, noise, k2mu0_Mu0VgammaDt, Msat, Temp, Alpha, &wg)
	} else {
		go settemperature__(Bth, noise, k2mu0_Mu0VgammaDt, Msat, Temp, Alpha, &wg)
	}
	wg.Done()
}

func settemperature__(Bth, noise *data.Slice, k2mu0_Mu0VgammaDt float64, Msat, Temp, Alpha MSlice, wg_ *sync.WaitGroup) {
	var Beff unsafe.Pointer
	var nois unsafe.Pointer
	var Msat_X unsafe.Pointer
	var Temp_X unsafe.Pointer
	var Alpha_X unsafe.Pointer
	if Bth != nil {
		Bth.Lock(0)
		defer Bth.Unlock(0)
		Beff = Bth.DevPtr(0)
	}
	if noise != nil {
		noise.RLock(0)
		defer noise.RUnlock(0)
		nois = noise.DevPtr(0)
	}
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
		Msat_X = Msat.DevPtr(0)
	}
	if Temp.GetSlicePtr() != nil {
		Temp.RLock()
		defer Temp.RUnlock()
		Temp_X = Temp.DevPtr(0)
	}
	if Alpha.GetSlicePtr() != nil {
		Alpha.RLock()
		defer Alpha.RUnlock()
		Alpha_X = Alpha.DevPtr(0)
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("settemperature failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := Bth.Len()
	cfg := make1DConf(N)

	event := k_settemperature2_async(Beff, nois, float32(k2mu0_Mu0VgammaDt),
		Msat_X, Msat.Mul(0),
		Temp_X, Temp.Mul(0),
		Alpha_X, Alpha.Mul(0),
		N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in settemperature: %+v \n", err)
	}
}
