package opencl

// Generation of Magnetic Force Microscopy images.

import (
	"fmt"

	data "github.com/seeder-research/uMagNUS/data"
	mag "github.com/seeder-research/uMagNUS/mag"
)

// Stores the necessary state to perform FFT-accelerated convolution
type MFMConvolution struct {
	size        [3]int         // 3D size of the input/output data
	kernSize    [3]int         // Size of kernel and logical FFT size.
	fftKernSize [3]int         //
	fftRBuf     *data.Slice    // FFT input buf for FFT, shares storage with fftCBuf.
	fftCBuf     *data.Slice    // FFT output buf, shares storage with fftRBuf
	gpuFFTKern  [3]*data.Slice // FFT kernel on device
	fwPlan      fft3DR2CPlan   // Forward FFT (1 component)
	bwPlan      fft3DC2RPlan   // Backward FFT (1 component)
	kern        [3]*data.Slice // Real-space kernel (host)
	mesh        *data.Mesh
}

func (c *MFMConvolution) Free() {
	if c == nil {
		return
	}
	c.size = [3]int{}
	c.kernSize = [3]int{}
	c.fftCBuf.Free() // shared with fftRbuf
	c.fftCBuf = nil
	c.fftRBuf = nil

	for j := 0; j < 3; j++ {
		c.gpuFFTKern[j].Free()
		c.gpuFFTKern[j] = nil
		c.kern[j] = nil
	}
	c.fwPlan.Free()
	c.bwPlan.Free()
}

func (c *MFMConvolution) init() {
	// init FFT plans
	padded := c.kernSize
	c.fwPlan = newFFT3DR2C(padded[X], padded[Y], padded[Z])
	c.bwPlan = newFFT3DC2R(padded[X], padded[Y], padded[Z])

	// init device buffers
	nc := fftR2COutputSizeFloats(c.kernSize)
	c.fftCBuf = NewSlice(1, nc)
	c.fftRBuf = NewSlice(1, c.kernSize)

	c.gpuFFTKern[X] = NewSlice(1, nc)
	c.gpuFFTKern[Y] = NewSlice(1, nc)
	c.gpuFFTKern[Z] = NewSlice(1, nc)

	c.initFFTKern3D()
}

func (c *MFMConvolution) initFFTKern3D() {
	c.fftKernSize = fftR2COutputSizeFloats(c.kernSize)

	for i := 0; i < 3; i++ {
		zero1_async(c.fftRBuf)
		data.Copy(c.fftRBuf, c.kern[i])
		if err := c.fwPlan.ExecAsync(c.fftRBuf, c.fftCBuf); err != nil {
			fmt.Printf("error enqueuing forward fft in initfftkern3d: %+v \n", err)
		}
		if err := c.fwPlan.Sync(); err != nil {
			fmt.Printf("Failed to wait for command queue to complete in initfftkern3d: %+v \n", err)
		}
		scale := 2 / float32(c.fwPlan.InputLen()) // ??
		zero1_async(c.gpuFFTKern[i])
		Madd2(c.gpuFFTKern[i], c.gpuFFTKern[i], c.fftCBuf, 0, scale)
	}
}

// store MFM image in output, based on magnetization in inp.
func (c *MFMConvolution) Exec(outp, inp, vol *data.Slice, Msat MSlice) {
	for i := 0; i < 3; i++ {
		zero1_async(c.fftRBuf)
		copyPadMul(c.fftRBuf, inp.Comp(i), vol, c.kernSize, c.size, Msat)
		if err := ClCmdQueue.Finish(); err != nil {
			fmt.Printf("Failed to wait for main command queue to finish in mfmconv exec: %+v \n", err)
		}
		if err := c.fwPlan.ExecAsync(c.fftRBuf, c.fftCBuf); err != nil {
			fmt.Printf("error enqueuing forward fft in mfmconv exec: %+v \n", err)
		}
		if err := c.fwPlan.Sync(); err != nil {
			fmt.Printf("error syncing queue for forward fft in mfmconv exec: %+v \n", err)
		}

		Nx, Ny := c.fftKernSize[X]/2, c.fftKernSize[Y] //   ??
		kernMulC_async(c.fftCBuf, c.gpuFFTKern[i], Nx, Ny)
		if err := ClCmdQueue.Finish(); err != nil {
			fmt.Printf("Failed to wait for main command queue to finish in mfmconv exec: %+v \n", err)
		}

		if err := c.bwPlan.ExecAsync(c.fftCBuf, c.fftRBuf); err != nil {
			fmt.Printf("error enqueuing backward fft in mfmconv exec: %+v \n", err)
		}
		if err := c.bwPlan.Sync(); err != nil {
			fmt.Printf("error syncing queue for backward fft in mfmconv exec: %+v \n", err)
		}
		copyUnPad(outp.Comp(i), c.fftRBuf, c.size, c.kernSize)
		if err := ClCmdQueue.Finish(); err != nil {
			fmt.Printf("Failed to wait for main command queue to finish in mfmconv exec: %+v \n", err)
		}
	}
}

func (c *MFMConvolution) Reinit(lift, tipsize float64, cachedir string) {
	c.kern = mag.MFMKernel(c.mesh, lift, tipsize, cachedir)
	c.initFFTKern3D()
}

// Initializes a convolution to evaluate the demag field for the given mesh geometry.
func NewMFM(mesh *data.Mesh, lift, tipsize float64, cachedir string) *MFMConvolution {
	k := mag.MFMKernel(mesh, lift, tipsize, cachedir)
	size := mesh.Size()
	c := new(MFMConvolution)
	c.size = size
	c.kern = k
	c.kernSize = k[X].Size()
	c.init()
	c.mesh = mesh
	return c
}
