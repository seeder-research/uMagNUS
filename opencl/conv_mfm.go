package opencl

// Generation of Magnetic Force Microscopy images.

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	mag "github.com/seeder-research/uMagNUS/mag"
)

// Stores the necessary state to perform FFT-accelerated convolution
type MFMConvolution struct {
	size        [3]int          // 3D size of the input/output data
	kernSize    [3]int          // Size of kernel and logical FFT size.
	fftKernSize [3]int          //
	fftRBuf     *data.Slice     // FFT input buf for FFT, shares storage with fftCBuf.
	fftCBuf     *data.Slice     // FFT output buf, shares storage with fftRBuf
	gpuFFTKern  [3]*data.Slice  // FFT kernel on device
	fwPlan      [3]fft3DR2CPlan // Forward FFT (1 component)
	bwPlan      [3]fft3DC2RPlan // Backward FFT (1 component)
	fwQueue     [3]*cl.CommandQueue
	bwQueue     [3]*cl.CommandQueue
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
		c.fwQueue[j].Release()
		c.bwQueue[j].Release()
		c.fwPlan[j].Free()
		c.bwPlan[j].Free()
	}
}

func (c *MFMConvolution) init(queue *cl.CommandQueue, events []*cl.Event) {
	// init FFT plans
	padded := c.kernSize
	c.fwPlan[X] = newFFT3DR2C(padded[X], padded[Y], padded[Z])
	c.fwPlan[Y] = newFFT3DR2C(padded[X], padded[Y], padded[Z])
	c.fwPlan[Z] = newFFT3DR2C(padded[X], padded[Y], padded[Z])
	c.bwPlan[X] = newFFT3DC2R(padded[X], padded[Y], padded[Z])
	c.bwPlan[Y] = newFFT3DC2R(padded[X], padded[Y], padded[Z])
	c.bwPlan[Z] = newFFT3DC2R(padded[X], padded[Y], padded[Z])
	c.fwQueue[X] = c.fwPlan[X].GetCommandQueue()
	c.fwQueue[Y] = c.fwPlan[Y].GetCommandQueue()
	c.fwQueue[Z] = c.fwPlan[Z].GetCommandQueue()
	c.bwQueue[X] = c.bwPlan[X].GetCommandQueue()
	c.bwQueue[Y] = c.bwPlan[Y].GetCommandQueue()
	c.bwQueue[Z] = c.bwPlan[Z].GetCommandQueue()

	// init device buffers
	nc := fftR2COutputSizeFloats(c.kernSize)
	c.fftCBuf = NewSlice(1, nc)
	c.fftRBuf = NewSlice(1, c.kernSize)

	c.gpuFFTKern[X] = NewSlice(1, nc)
	c.gpuFFTKern[Y] = NewSlice(1, nc)
	c.gpuFFTKern[Z] = NewSlice(1, nc)

	c.initFFTKern3D(queue, events)
}

func (c *MFMConvolution) initFFTKern3D(queue *cl.CommandQueue, events []*cl.Event) {
	c.fftKernSize = fftR2COutputSizeFloats(c.kernSize)

	SyncQueues([]*cl.CommandQueue{ClCmdQueue[1], ClCmdQueue[2], ClCmdQueue[3]}, []*cl.CommandQueue{queue})
	for i := 0; i < 3; i++ {
		zero1_async(c.fftRBuf, ClCmdQueue[1+i], events)
		data.Copy(c.fftRBuf, c.kern[i])
		SyncQueues([]*cl.CommandQueue{c.fwQueue[i]}, []*cl.CommandQueue{ClCmdQueue[1+i]})
		err := c.fwPlan[i].ExecAsync(c.fftRBuf, c.fftCBuf, ClCmdQueue[1+i], nil)
		if err != nil {
			fmt.Printf("error enqueuing forward fft in initfftkern3d: %+v \n", err)
		}
		scale := 2 / float32(c.fwPlan[i].InputLen()) // ??
		SyncQueues([]*cl.CommandQueue{ClCmdQueue[1+i]}, []*cl.CommandQueue{c.fwQueue[i]})
		zero1_async(c.gpuFFTKern[i], ClCmdQueue[1+i], nil)
		SyncQueues([]*cl.CommandQueue{ClCmdQueue[1+i]}, []*cl.CommandQueue{queue})
		Madd2(c.gpuFFTKern[i], c.gpuFFTKern[i], c.fftCBuf, 0, scale, []*cl.CommandQueue{ClCmdQueue[1+i]}, nil)
	}
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{ClCmdQueue[1], ClCmdQueue[2], ClCmdQueue[3]})
}

// store MFM image in output, based on magnetization in inp.
func (c *MFMConvolution) Exec(outp, inp, vol *data.Slice, Msat MSlice, queue *cl.CommandQueue, events []*cl.Event) {
	SyncQueues([]*cl.CommandQueue{ClCmdQueue[1], ClCmdQueue[2], ClCmdQueue[3]}, []*cl.CommandQueue{queue})
	for i := 0; i < 3; i++ {
		zero1_async(c.fftRBuf, ClCmdQueue[1+i], events)
		copyPadMul(c.fftRBuf, inp.Comp(i), vol, c.kernSize, c.size, Msat, ClCmdQueue[1+i], nil)
		SyncQueues([]*cl.CommandQueue{c.fwQueue[i]}, []*cl.CommandQueue{ClCmdQueue[1+i]})
		var err error
		if err = c.fwPlan[i].ExecAsync(c.fftRBuf, c.fftCBuf, ClCmdQueue[1+i], nil); err != nil {
			fmt.Printf("error enqueuing forward fft in mfmconv exec: %+v \n", err)
		}

		Nx, Ny := c.fftKernSize[X]/2, c.fftKernSize[Y] //   ??
		SyncQueues([]*cl.CommandQueue{ClCmdQueue[1+i]}, []*cl.CommandQueue{c.fwQueue[i]})
		kernMulC_async(c.fftCBuf, c.gpuFFTKern[i], Nx, Ny, ClCmdQueue[1+i], nil)

		SyncQueues([]*cl.CommandQueue{c.bwQueue[i]}, []*cl.CommandQueue{ClCmdQueue[1+i]})
		if err = c.bwPlan[i].ExecAsync(c.fftCBuf, c.fftRBuf, ClCmdQueue[1+i], nil); err != nil {
			fmt.Printf("error enqueuing backward fft in mfmconv exec: %+v \n", err)
		}
		SyncQueues([]*cl.CommandQueue{ClCmdQueue[1+i]}, []*cl.CommandQueue{c.bwQueue[i]})
		copyUnPad(outp.Comp(i), c.fftRBuf, c.size, c.kernSize, ClCmdQueue[1+i], nil)
	}
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{ClCmdQueue[1], ClCmdQueue[2], ClCmdQueue[3]})
}

func (c *MFMConvolution) Reinit(lift, tipsize float64, cachedir string, queue *cl.CommandQueue, events []*cl.Event) {
	c.kern = mag.MFMKernel(c.mesh, lift, tipsize, cachedir)
	c.initFFTKern3D(queue, events)
}

// Initializes a convolution to evaluate the demag field for the given mesh geometry.
func NewMFM(mesh *data.Mesh, lift, tipsize float64, cachedir string, queue *cl.CommandQueue, events []*cl.Event) *MFMConvolution {
	k := mag.MFMKernel(mesh, lift, tipsize, cachedir)
	size := mesh.Size()
	c := new(MFMConvolution)
	c.size = size
	c.kern = k
	c.kernSize = k[X].Size()
	c.init(queue, events)
	c.mesh = mesh
	return c
}
