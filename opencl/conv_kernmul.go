package opencl

// Kernel multiplication for purely real kernel, symmetric around Y axis (apart from first row).
// Launch configs range over all complex elements of fft input. This could be optimized: range only over kernel.

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// kernel multiplication for 3D demag convolution, exploiting full kernel symmetry.
func kernMulRSymm3D_async(fftM [3]*data.Slice, Kxx, Kyy, Kzz, Kyz, Kxz, Kxy *data.Slice, Nx, Ny, Nz int, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(fftM[X].NComp() == 1 && Kxx.NComp() == 1)
	cfg := make3DConf([3]int{Nx, Ny, Nz})

	// Launch kernel
	event := k_kernmulRSymm3D_async(fftM[X].DevPtr(0), fftM[Y].DevPtr(0), fftM[Z].DevPtr(0),
		Kxx.DevPtr(0), Kyy.DevPtr(0), Kzz.DevPtr(0), Kyz.DevPtr(0), Kxz.DevPtr(0), Kxy.DevPtr(0),
		Nx, Ny, Nz, cfg, ewl, q)

	fftM[X].SetEvent(0, event)
	fftM[Y].SetEvent(0, event)
	fftM[Z].SetEvent(0, event)

	glist := []GSlice{Kxx, Kyy, Kzz, Kyz, Kxz, Kxy}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in kernmulrsymm3d_async: %+v \n", err)
		}
	}

	return
}

// kernel multiplication for 2D demag convolution on X and Y, exploiting full kernel symmetry.
func kernMulRSymm2Dxy_async(fftMx, fftMy, Kxx, Kyy, Kxy *data.Slice, Nx, Ny int, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(fftMy.NComp() == 1 && Kxx.NComp() == 1)
	cfg := make3DConf([3]int{Nx, Ny, 1})


	event := k_kernmulRSymm2Dxy_async(fftMx.DevPtr(0), fftMy.DevPtr(0),
		Kxx.DevPtr(0), Kyy.DevPtr(0), Kxy.DevPtr(0),
		Nx, Ny, cfg, ewl,
		q)

	fftMx.SetEvent(0, event)
	fftMy.SetEvent(0, event)

	glist := []GSlice{Kxx, Kyy, Kxy}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in kernmulrsymm2dxy_async: %+v \n", err)
		}
	}

	return
}

// kernel multiplication for 2D demag convolution on Z, exploiting full kernel symmetry.
func kernMulRSymm2Dz_async(fftMz, Kzz *data.Slice, Nx, Ny int, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(fftMz.NComp() == 1 && Kzz.NComp() == 1)
	cfg := make3DConf([3]int{Nx, Ny, 1})

	event := k_kernmulRSymm2Dz_async(fftMz.DevPtr(0), Kzz.DevPtr(0), Nx, Ny, cfg, ewl, q)

	fftMz.SetEvent(0, event)

	glist := []GSlice{Kzz}
	InsertEventIntoGSlices(event, glist)

	if Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in kernmulrsymm2dz_async: %+v \n", err)
		}
	}

	return
}

// kernel multiplication for general 1D convolution. Does not assume any symmetry.
// Used for MFM images.
func kernMulC_async(fftM, K *data.Slice, Nx, Ny int, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(fftM.NComp() == 1 && K.NComp() == 1)
	cfg := make3DConf([3]int{Nx, Ny, 1})

	event := k_kernmulC_async(fftM.DevPtr(0), K.DevPtr(0), Nx, Ny, cfg, ewl, q)

	fftM.SetEvent(0, event)

	glist := []GSlice{K}
	InsertEventIntoGSlices(event, glist)

	if Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in kernmulC_async: %+v \n", err)
		}
	}

	return
}
