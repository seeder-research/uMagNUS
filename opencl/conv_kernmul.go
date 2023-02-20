package opencl

// Kernel multiplication for purely real kernel, symmetric around Y axis (apart from first row).
// Launch configs range over all complex elements of fft input. This could be optimized: range only over kernel.

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// kernel multiplication for 3D demag convolution, exploiting full kernel symmetry.
func kernMulRSymm3D_async(fftM [3]*data.Slice, Kxx, Kyy, Kzz, Kyz, Kxz, Kxy *data.Slice, Nx, Ny, Nz int) {
	util.Argument(fftM[X].NComp() == 1 && Kxx.NComp() == 1)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		kernmulrsymm3d_async__(fftM, Kxx, Kyy, Kzz, Kyz, Kxz, Kxy, Nx, Ny, Nz, &wg)
	} else {
		go kernmulrsymm3d_async__(fftM, Kxx, Kyy, Kzz, Kyz, Kxz, Kxy, Nx, Ny, Nz, &wg)
	}
	wg.Wait()
}

func kernmulrsymm3d_async__(fftM [3]*data.Slice, Kxx, Kyy, Kzz, Kyz, Kxz, Kxy *data.Slice, Nx, Ny, Nz int, wg_ *sync.WaitGroup) {
	fftM[X].Lock(0)
	fftM[Y].Lock(0)
	fftM[Z].Lock(0)
	defer fftM[X].Unlock(0)
	defer fftM[Y].Unlock(0)
	defer fftM[Z].Unlock(0)
	Kxx.RLock(0)
	Kyy.RLock(0)
	Kzz.RLock(0)
	Kxy.RLock(0)
	Kxz.RLock(0)
	Kyz.RLock(0)
	defer Kxx.RUnlock(0)
	defer Kyy.RUnlock(0)
	defer Kzz.RUnlock(0)
	defer Kxy.RUnlock(0)
	defer Kxz.RUnlock(0)
	defer Kyz.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("kernmulrsymm3d_async failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	cfg := make3DConf([3]int{Nx, Ny, Nz})

	event := k_kernmulRSymm3D_async(fftM[X].DevPtr(0), fftM[Y].DevPtr(0), fftM[Z].DevPtr(0),
		Kxx.DevPtr(0), Kyy.DevPtr(0), Kzz.DevPtr(0), Kyz.DevPtr(0), Kxz.DevPtr(0), Kxy.DevPtr(0),
		Nx, Ny, Nz, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in kernmulrsymm3d_async: %+v \n", err)
	}
}

// kernel multiplication for 2D demag convolution on X and Y, exploiting full kernel symmetry.
func kernMulRSymm2Dxy_async(fftMx, fftMy, Kxx, Kyy, Kxy *data.Slice, Nx, Ny int) {
	util.Argument(fftMy.NComp() == 1 && Kxx.NComp() == 1)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		kernmulrsymm2dxy_async__(fftMx, fftMy, Kxx, Kyy, Kxy, Nx, Ny, &wg)
	} else {
		go kernmulrsymm2dxy_async__(fftMx, fftMy, Kxx, Kyy, Kxy, Nx, Ny, &wg)
	}
	wg.Wait()
}

func kernmulrsymm2dxy_async__(fftMx, fftMy, Kxx, Kyy, Kxy *data.Slice, Nx, Ny int, wg_ *sync.WaitGroup) {
	fftMx.Lock(0)
	fftMy.Lock(0)
	defer fftMx.Unlock(0)
	defer fftMy.Unlock(0)
	Kxx.RLock(0)
	Kyy.RLock(0)
	Kxy.RLock(0)
	defer Kxx.RUnlock(0)
	defer Kyy.RUnlock(0)
	defer Kxy.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("kernmulrsymm2dxy_async failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	cfg := make3DConf([3]int{Nx, Ny, 1})

	event := k_kernmulRSymm2Dxy_async(fftMx.DevPtr(0), fftMy.DevPtr(0),
		Kxx.DevPtr(0), Kyy.DevPtr(0), Kxy.DevPtr(0),
		Nx, Ny, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in kernmulrsymm2dxy_async: %+v \n", err)
	}
}

// kernel multiplication for 2D demag convolution on Z, exploiting full kernel symmetry.
func kernMulRSymm2Dz_async(fftMz, Kzz *data.Slice, Nx, Ny int) {
	util.Argument(fftMz.NComp() == 1 && Kzz.NComp() == 1)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		kernmulrsymm2dz_async__(fftMz, Kzz, Nx, Ny, &wg)
	} else {
		go kernmulrsymm2dz_async__(fftMz, Kzz, Nx, Ny, &wg)
	}
	wg.Wait()
}

func kernmulrsymm2dz_async__(fftMz, Kzz *data.Slice, Nx, Ny int, wg_ *sync.WaitGroup) {
	fftMz.Lock(0)
	defer fftMz.Unlock(0)
	Kzz.RLock(0)
	defer Kzz.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("kernmulrsymm2dz_async failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	cfg := make3DConf([3]int{Nx, Ny, 1})

	event := k_kernmulRSymm2Dz_async(fftMz.DevPtr(0), Kzz.DevPtr(0), Nx, Ny, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in kernmulrsymm2dz_async: %+v \n", err)
	}
}

// kernel multiplication for general 1D convolution. Does not assume any symmetry.
// Used for MFM images.
func kernMulC_async(fftM, K *data.Slice, Nx, Ny int) {
	util.Argument(fftM.NComp() == 1 && K.NComp() == 1)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		kernmulc_async__(fftM, K, Nx, Ny, &wg)
	} else {
		go kernmulc_async__(fftM, K, Nx, Ny, &wg)
	}
	wg.Wait()
}

func kernmulc_async__(fftM, K *data.Slice, Nx, Ny int, wg_ *sync.WaitGroup) {
	fftM.Lock(0)
	defer fftM.Unlock(0)
	K.RLock(0)
	defer K.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("kernmulrsymm2dz_async failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	cfg := make3DConf([3]int{Nx, Ny, 1})

	event := k_kernmulC_async(fftM.DevPtr(0), K.DevPtr(0), Nx, Ny, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in kernmulC_async: %+v \n", err)
	}
}
