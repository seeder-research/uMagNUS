package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Stores the necessary state to perform FFT-accelerated convolution
// with magnetostatic kernel (or other kernel of same symmetry).
type DemagConvolution struct {
	inputSize        [3]int            // 3D size of the input/output data
	realKernSize     [3]int            // Size of kernel and logical FFT size.
	fftKernLogicSize [3]int            // logic size FFTed kernel, real parts only, we store less
	fftRBuf          [3]*data.Slice    // FFT input buf; 2D: Z shares storage with X.
	fftCBuf          [3]*data.Slice    // FFT output buf; 2D: Z shares storage with X.
	kern             [3][3]*data.Slice // FFT kernel on device
	fwPlan           [3]fft3DR2CPlan   // Forward FFT (1 for each component)
	bwPlan           [3]fft3DC2RPlan   // Backward FFT (1 for each component)
	fwQueue          [3]*cl.CommandQueue
	bwQueue          [3]*cl.CommandQueue
}

// Initializes a convolution to evaluate the demag field for the given mesh geometry.
// Sanity-checked if test == true (slow-ish for large meshes).
func NewDemag(inputSize, PBC [3]int, kernel [3][3]*data.Slice, test bool, queue *cl.CommandQueue, events []*cl.Event) *DemagConvolution {
	c := new(DemagConvolution)
	c.inputSize = inputSize
	c.realKernSize = kernel[X][X].Size()
	c.init(kernel, queue, events)
	if test {
		testConvolution(c, PBC, kernel, queue, events)
	}
	return c
}

// Calculate the demag field of m * vol * Bsat, store result in B.
//
//	m:    magnetization normalized to unit length
//	vol:  unitless mask used to scale m's length, may be nil
//	Bsat: saturation magnetization in Tesla
//	B:    resulting demag field, in Tesla
func (c *DemagConvolution) Exec(B, m, vol *data.Slice, Msat MSlice, queue *cl.CommandQueue, events []*cl.Event) {
	util.Argument(B.Size() == c.inputSize && m.Size() == c.inputSize)
	if c.is2D() {
		c.exec2D(B, m, vol, Msat, queue, events)
	} else {
		c.exec3D(B, m, vol, Msat, queue, events)
	}
}

func (c *DemagConvolution) exec3D(outp, inp, vol *data.Slice, Msat MSlice, queue *cl.CommandQueue, events []*cl.Event) {
	if Synchronous {
		if err := WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in demagconvolution.exec3d: %+v \n", err)
		}
	}

	if events != nil {
		if err := cl.WaitForEvents(events); err != nil {
			fmt.Printf("failed to WaitForEvents in beginning of exec2d: %+v \n", err)
		}
	}

	for i := 0; i < 3; i++ { // FW FFT
		c.fwFFT(i, inp, vol, Msat, queue, nil)
	}
	// Synchronize main queue to fwFFT queues
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{c.fwQueue[X], c.fwQueue[Y], c.fwQueue[Z]})

	// kern mul
	kernMulRSymm3D_async(c.fftCBuf,
		c.kern[X][X], c.kern[Y][Y], c.kern[Z][Z],
		c.kern[Y][Z], c.kern[X][Z], c.kern[X][Y],
		c.fftKernLogicSize[X], c.fftKernLogicSize[Y], c.fftKernLogicSize[Z],
		queue, events)

	// Synchronize bwFFT queues to main queue
	SyncQueues([]*cl.CommandQueue{c.fwQueue[X], c.fwQueue[Y], c.fwQueue[Z]}, []*cl.CommandQueue{queue})

	for i := 0; i < 3; i++ { // BW FFT
		c.bwFFT(i, outp, queue, nil)
	}

	// Synchronize main queue to bwFFT queues
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{c.bwQueue[X], c.bwQueue[Y], c.bwQueue[Z]})
}

func (c *DemagConvolution) exec2D(outp, inp, vol *data.Slice, Msat MSlice, queue *cl.CommandQueue, events []*cl.Event) {
	if Synchronous {
		if err := WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in demagconvolution.exec2d: %+v \n", err)
		}
	}

	if events != nil {
		if err := cl.WaitForEvents(events); err != nil {
			fmt.Printf("failed to WaitForEvents in beginning of exec2d: %+v \n", err)
		}
	}

	// Convolution is separated into
	// a 1D convolution for z and a 2D convolution for xy.
	// So only 2 FFT buffers are needed at the same time.
	Nx, Ny := c.fftKernLogicSize[X], c.fftKernLogicSize[Y]

	// Z
	c.fwFFT(Z, inp, vol, Msat, queue, nil)

	// Sync to fwFFT queue
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{c.fwQueue[Z]})

	kernMulRSymm2Dz_async(c.fftCBuf[Z], c.kern[Z][Z], Nx, Ny, queue, nil)

	// Sync bwFFT queue
	SyncQueues([]*cl.CommandQueue{c.bwQueue[Z]}, []*cl.CommandQueue{queue})

	c.bwFFT(Z, outp, queue, nil)

	// XY
	c.fwFFT(X, inp, vol, Msat, queue, nil)
	c.fwFFT(Y, inp, vol, Msat, queue, nil)
	// Synchronize main queue to fwFFT queues
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{c.fwQueue[X], c.fwQueue[Y]})

	kernMulRSymm2Dxy_async(c.fftCBuf[X], c.fftCBuf[Y],
		c.kern[X][X], c.kern[Y][Y], c.kern[X][Y], Nx, Ny, queue, nil)

	// Synchronize bwFFT queues to main queue
	SyncQueues([]*cl.CommandQueue{c.bwQueue[X], c.bwQueue[Y]}, []*cl.CommandQueue{queue})
	c.bwFFT(X, outp, queue, nil)
	c.bwFFT(Y, outp, queue, nil)

	// Synchronize main queue to bwFFT queues
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{c.bwQueue[X], c.bwQueue[Y], c.bwQueue[Z]})
}

func (c *DemagConvolution) is2D() bool {
	return c.inputSize[Z] == 1
}

// zero 1-component slice
func zero1_async(dst *data.Slice, queue *cl.CommandQueue, events []*cl.Event) {
	val := float32(0.0)
	if dst == nil {
		panic("ERROR (zero1_async): dst pointer cannot be nil")
	}
	if Synchronous {
		if err := WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in zero1_async: %+v \n", err)
		}
	}

	event, err := queue.EnqueueFillBuffer((*cl.MemObject)(dst.DevPtr(0)), unsafe.Pointer(&val), SIZEOF_FLOAT32, 0, dst.Len()*SIZEOF_FLOAT32, events)
	if err != nil {
		fmt.Printf("EnqueueFillBuffer failed: %+v \n", err)
	}
	if Synchronous {
		if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in zero1_async: %+v \n", err)
		}
	}
}

// forward FFT component i
func (c *DemagConvolution) fwFFT(i int, inp, vol *data.Slice, Msat MSlice, queue *cl.CommandQueue, events []*cl.Event) {
	zero1_async(c.fftRBuf[i], queue, events)
	in := inp.Comp(i)
	copyPadMul(c.fftRBuf[i], in, vol, c.realKernSize, c.inputSize, Msat, queue, nil)
	// Sync fwFFT queue to main (sequential) queue and execute fwfft
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{c.fwQueue[i]})
	if err := c.fwPlan[i].ExecAsync(c.fftRBuf[i], c.fftCBuf[i], queue, nil); err != nil {
		fmt.Printf("Error enqueuing forward fft: %+v \n", err)
	}
}

// backward FFT component i
func (c *DemagConvolution) bwFFT(i int, outp *data.Slice, queue *cl.CommandQueue, events []*cl.Event) {
	if events != nil {
		if err := cl.WaitForEvents(events); err != nil {
			fmt.Printf("WaitForEvents failed in beginning of bwFFT")
		}
	}
	if err := c.bwPlan[i].ExecAsync(c.fftCBuf[i], c.fftRBuf[i], queue, nil); err != nil {
		fmt.Printf("Error enqueuing backward fft: %+v", err)
	}
	out := outp.Comp(i)
	// Sync main (sequential) queue to bwFFT queue and execute bwfft
	SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{c.bwQueue[i]})
	copyUnPad(out, c.fftRBuf[i], c.inputSize, c.realKernSize, queue, nil)
}

func (c *DemagConvolution) init(realKern [3][3]*data.Slice, queue *cl.CommandQueue, events []*cl.Event) {
	// init device buffers
	// 2D re-uses fftBuf[X] as fftBuf[Z], 3D needs all 3 fftBufs.
	nc := fftR2COutputSizeFloats(c.realKernSize)
	c.fftCBuf[X] = NewSlice(1, nc)
	c.fftCBuf[Y] = NewSlice(1, nc)
	if c.is2D() {
		c.fftCBuf[Z] = c.fftCBuf[X]
	} else {
		c.fftCBuf[Z] = NewSlice(1, nc)
	}

	c.fftRBuf[X] = NewSlice(1, c.realKernSize)
	c.fftRBuf[Y] = NewSlice(1, c.realKernSize)
	if c.is2D() {
		c.fftRBuf[Z] = c.fftRBuf[X]
	} else {
		c.fftRBuf[Z] = NewSlice(1, c.realKernSize)
	}

	// init FFT plans
	c.fwPlan[X] = newFFT3DR2C(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.fwPlan[Y] = newFFT3DR2C(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.fwPlan[Z] = newFFT3DR2C(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.bwPlan[X] = newFFT3DC2R(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.bwPlan[Y] = newFFT3DC2R(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.bwPlan[Z] = newFFT3DC2R(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.fwQueue[X] = c.fwPlan[X].GetCommandQueue()
	c.fwQueue[Y] = c.fwPlan[Y].GetCommandQueue()
	c.fwQueue[Z] = c.fwPlan[Z].GetCommandQueue()
	c.bwQueue[X] = c.bwPlan[X].GetCommandQueue()
	c.bwQueue[Y] = c.bwPlan[Y].GetCommandQueue()
	c.bwQueue[Z] = c.bwPlan[Z].GetCommandQueue()

	// init FFT kernel

	// logic size of FFT(kernel): store real parts only
	c.fftKernLogicSize = fftR2COutputSizeFloats(c.realKernSize)
	util.Assert(c.fftKernLogicSize[X]%2 == 0)
	c.fftKernLogicSize[X] /= 2

	// physical size of FFT(kernel): store only non-redundant part exploiting Y, Z mirror symmetry
	// X mirror symmetry already exploited: FFT(kernel) is purely real.
	physKSize := [3]int{c.fftKernLogicSize[X], c.fftKernLogicSize[Y]/2 + 1, c.fftKernLogicSize[Z]/2 + 1}

	output := c.fftCBuf[0]
	input := c.fftRBuf[0]
	fftKern := data.NewSlice(1, physKSize)
	kfull := data.NewSlice(1, output.Size()) // not yet exploiting symmetry
	kfulls := kfull.Scalars()
	kCSize := physKSize
	kCSize[X] *= 2                     // size of kernel after removing Y,Z redundant parts, but still complex
	kCmplx := data.NewSlice(1, kCSize) // not yet exploiting X symmetry
	kc := kCmplx.Scalars()

	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ { // upper triangular part
			if realKern[i][j] != nil { // ignore 0's
				// FW FFT
				data.Copy(input, realKern[i][j])
				SyncQueues([]*cl.CommandQueue{c.fwQueue[j]}, []*cl.CommandQueue{queue})
				err := c.fwPlan[j].ExecAsync(input, output, queue, nil) // FW FFT
				if err != nil {
					fmt.Printf("error enqueuing forward fft in init: %+v \n ", err)
				}
				// Sync main (sequential) queue to fwFFT queue
				SyncQueues([]*cl.CommandQueue{queue}, []*cl.CommandQueue{c.fwQueue[j]})
				data.Copy(kfull, output)
				// Wait for FFT to complete and copyback to complete
				if err = queue.Finish(); err != nil {
					fmt.Printf("error waiting main queue to finish after fft copyback in init: %+v \n ", err)
				}

				// extract non-redundant part (Y,Z symmetry)
				for iz := 0; iz < kCSize[Z]; iz++ {
					for iy := 0; iy < kCSize[Y]; iy++ {
						for ix := 0; ix < kCSize[X]; ix++ {
							kc[iz][iy][ix] = kfulls[iz][iy][ix]
						}
					}
				}

				// extract real parts (X symmetry)
				scaleRealParts(fftKern, kCmplx, 1/float32(c.fwPlan[j].InputLen()))
				c.kern[i][j] = GPUCopy(fftKern)
			}
		}
	}
}

func (c *DemagConvolution) Free() {
	if c == nil {
		return
	}
	c.inputSize = [3]int{}
	c.realKernSize = [3]int{}
	for i := 0; i < 3; i++ {
		c.fftCBuf[i].Free()
		c.fftRBuf[i].Free()
		c.fftCBuf[i] = nil
		c.fftRBuf[i] = nil

		for j := 0; j < 3; j++ {
			c.kern[i][j].Free()
			c.kern[i][j] = nil
		}
		c.fwQueue[i].Release()
		c.bwQueue[i].Release()
		c.fwPlan[i].Free()
		c.bwPlan[i].Free()
	}
}
