package opencl64

import (
	"fmt"
	"unsafe"

	data "github.com/seeder-research/uMagNUS/data64"
	"github.com/seeder-research/uMagNUS/cl"
	"github.com/seeder-research/uMagNUS/util"
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
	fwPlan           fft3DR2CPlan      // Forward FFT (1 component)
	bwPlan           fft3DC2RPlan      // Backward FFT (1 component)
}

// Initializes a convolution to evaluate the demag field for the given mesh geometry.
// Sanity-checked if test == true (slow-ish for large meshes).
func NewDemag(inputSize, PBC [3]int, kernel [3][3]*data.Slice, test bool) *DemagConvolution {
	c := new(DemagConvolution)
	c.inputSize = inputSize
	c.realKernSize = kernel[X][X].Size()
	c.init(kernel)
	if test {
		testConvolution(c, PBC, kernel)
	}
	return c
}

// Calculate the demag field of m * vol * Bsat, store result in B.
// 	m:    magnetization normalized to unit length
// 	vol:  unitless mask used to scale m's length, may be nil
// 	Bsat: saturation magnetization in Tesla
// 	B:    resulting demag field, in Tesla
func (c *DemagConvolution) Exec(B, m, vol *data.Slice, Msat MSlice) {
	util.Argument(B.Size() == c.inputSize && m.Size() == c.inputSize)
	if c.is2D() {
		c.exec2D(B, m, vol, Msat)
	} else {
		c.exec3D(B, m, vol, Msat)
	}
}

func (c *DemagConvolution) exec3D(outp, inp, vol *data.Slice, Msat MSlice) {
	//	events := make([][]*cl.Event, 3)
	//	totalLen := 0
	for i := 0; i < 3; i++ { // FW FFT
		//		events[i] = c.fwFFT(i, inp, vol, Msat)
		c.fwFFT(i, inp, vol, Msat)
		//		totalLen += len(events[i])
	}
	//	fullEventList := make([]*cl.Event, totalLen)
	//	for i := 0; i < 3; i++ {
	//		for _, eid := range events[i] {
	//			totalLen--
	//			fullEventList[totalLen] = eid
	//		}
	//	}
	//	err := cl.WaitForEvents(fullEventList)
	//	if err != nil {
	//		fmt.Printf("error waiting for forward fft to finish in exec3d:  %+v \n", err)
	//	}
	// kern mul
	kernMulRSymm3D_async(c.fftCBuf,
		c.kern[X][X], c.kern[Y][Y], c.kern[Z][Z],
		c.kern[Y][Z], c.kern[X][Z], c.kern[X][Y],
		c.fftKernLogicSize[X], c.fftKernLogicSize[Y], c.fftKernLogicSize[Z])

	for i := 0; i < 3; i++ { // BW FFT
		c.bwFFT(i, outp)
	}
}

func (c *DemagConvolution) exec2D(outp, inp, vol *data.Slice, Msat MSlice) {
	// Convolution is separated into
	// a 1D convolution for z and a 2D convolution for xy.
	// So only 2 FFT buffers are needed at the same time.
	Nx, Ny := c.fftKernLogicSize[X], c.fftKernLogicSize[Y]

	// Z
	//	event := c.fwFFT(Z, inp, vol, Msat)
	c.fwFFT(Z, inp, vol, Msat)
	//	err := cl.WaitForEvents(event)
	//	if err != nil {
	//		fmt.Printf("error waiting for forward fft to end in exec2d: %+v \n", err)
	//	}
	kernMulRSymm2Dz_async(c.fftCBuf[Z], c.kern[Z][Z], Nx, Ny)
	c.bwFFT(Z, outp)

	// XY
	//	event = c.fwFFT(X, inp, vol, Msat)
	//	event0 := c.fwFFT(Y, inp, vol, Msat)
	c.fwFFT(X, inp, vol, Msat)
	c.fwFFT(Y, inp, vol, Msat)
	//	offset := len(event)
	//	totalLen := offset + len(event0)
	//	fullEventList := make([]*cl.Event, totalLen)
	//	for id, eid := range event {
	//		fullEventList[id] = eid
	//	}
	//	for id, eid := range event0 {
	//		fullEventList[offset+id] = eid
	//	}
	//	err = cl.WaitForEvents(fullEventList)
	//	if err != nil {
	//		fmt.Printf("error waiting for second and third forward fft to end in exec2d: %+v \n", err)
	//	}
	kernMulRSymm2Dxy_async(c.fftCBuf[X], c.fftCBuf[Y],
		c.kern[X][X], c.kern[Y][Y], c.kern[X][Y], Nx, Ny)
	c.bwFFT(X, outp)
	c.bwFFT(Y, outp)
}

func (c *DemagConvolution) is2D() bool {
	return c.inputSize[Z] == 1
}

// zero 1-component slice
func zero1_async(dst *data.Slice) {
	val := float64(0.0)
	event, err := ClCmdQueue.EnqueueFillBuffer((*cl.MemObject)(dst.DevPtr(0)), unsafe.Pointer(&val), SIZEOF_FLOAT32, 0, dst.Len()*SIZEOF_FLOAT32, [](*cl.Event){dst.GetEvent(0)})
	dst.SetEvent(0, event)
	if err != nil {
		fmt.Printf("EnqueueFillBuffer failed: %+v \n", err)
	}
}

// forward FFT component i
func (c *DemagConvolution) fwFFT(i int, inp, vol *data.Slice, Msat MSlice) {
	//[]*cl.Event {
	zero1_async(c.fftRBuf[i])
	in := inp.Comp(i)
	copyPadMul(c.fftRBuf[i], in, vol, c.realKernSize, c.inputSize, Msat)
	//	event, err := c.fwPlan.ExecAsync(c.fftRBuf[i], c.fftCBuf[i])
	err := c.fwPlan.ExecAsync(c.fftRBuf[i], c.fftCBuf[i])
	if err != nil {
		fmt.Printf("Error enqueuing forward fft: %+v \n", err)
	}
	//	return event
}

// backward FFT component i
func (c *DemagConvolution) bwFFT(i int, outp *data.Slice) {
	//	event, err := c.bwPlan.ExecAsync(c.fftCBuf[i], c.fftRBuf[i])
	err := c.bwPlan.ExecAsync(c.fftCBuf[i], c.fftRBuf[i])
	if err != nil {
		fmt.Printf("Error enqueuing backward fft: %+v", err)
	}
	//	err = cl.WaitForEvents(event)
	//	if err != nil {
	//		fmt.Printf("Error waiting for backward fft t end: %+v \n ", err)
	//	}
	out := outp.Comp(i)
	copyUnPad(out, c.fftRBuf[i], c.inputSize, c.realKernSize)
}

func (c *DemagConvolution) init(realKern [3][3]*data.Slice) {
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
	c.fwPlan = newFFT3DR2C(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.bwPlan = newFFT3DC2R(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])

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
				//				event, err := c.fwPlan.ExecAsync(input, output)
				err := c.fwPlan.ExecAsync(input, output)
				if err != nil {
					fmt.Printf("error enqueuing forward fft in init: %+v \n ", err)
				}
				//				err = cl.WaitForEvents(event)
				//				if err != nil {
				//					fmt.Printf("error waiting for forward fft to end in init: %+v \n ", err)
				//				}
				data.Copy(kfull, output)

				// extract non-redundant part (Y,Z symmetry)
				for iz := 0; iz < kCSize[Z]; iz++ {
					for iy := 0; iy < kCSize[Y]; iy++ {
						for ix := 0; ix < kCSize[X]; ix++ {
							kc[iz][iy][ix] = kfulls[iz][iy][ix]
						}
					}
				}

				// extract real parts (X symmetry)
				scaleRealParts(fftKern, kCmplx, 1/float64(c.fwPlan.InputLen()))
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
		c.fwPlan.Free()
		c.bwPlan.Free()
	}
}
