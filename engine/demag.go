package engine

// Calculation of magnetostatic field

import (
	"log"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	mag "github.com/seeder-research/uMagNUS/mag"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

// Demag variables
var (
	Msat        = NewScalarParam("Msat", "A/m", "Saturation magnetization")
	M_full      = NewVectorField("m_full", "A/m", "Unnormalized magnetization", SetMFull)
	B_demag     = NewVectorField("B_demag", "T", "Magnetostatic field", SetDemagField)
	Edens_demag = NewScalarField("Edens_demag", "J/m3", "Magnetostatic energy density", AddEdens_demag)
	E_demag     = NewScalarValue("E_demag", "J", "Magnetostatic energy", GetDemagEnergy)

	EnableDemag  = true // enable/disable global demag field
	NoDemagSpins = NewScalarParam("NoDemagSpins", "", "Disable magnetostatic interaction per region (default=0, set to 1 to disable). "+
		"E.g.: NoDemagSpins.SetRegion(5, 1) disables the magnetostatic interaction in region 5.")
	conv_             *opencl.DemagConvolution // does the heavy lifting
	DemagAccuracy     = 6.0                    // Demag accuracy (divide cubes in at most N^3 points)
	EnableNewellDemag = false                  // enable/disable global demag field calculated using Newell formulation
	asymptotic_radius = 32                     // Radius (in number of cells) beyond which demag calculations fall back to far-field approximation
	zero_self_demag   = 0                      // Include/exclude self-demag
)

var AddEdens_demag = makeEdensAdder(&B_demag, -0.5)

func init() {

	DeclVar("EnableDemag", &EnableDemag, "Enables/disables demag (default=true)")
	DeclVar("EnableNewellDemag", &EnableNewellDemag, "Enables/disables demag using Newell's formulation (default=false)")
	DeclVar("DemagAccuracy", &DemagAccuracy, "Controls accuracy of demag kernel")
	registerEnergy(GetDemagEnergy, AddEdens_demag)
}

// Sets dst to the current demag field
func SetDemagField(dst *data.Slice) {
	seqQueue := opencl.ClCmdQueue[0]
	if EnableDemag || EnableNewellDemag {
		msat := Msat.MSlice()
		defer msat.Recycle()
		// sync in the beginning
		if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
			log.Printf("error waiting for data copy in setdemagfield 0: %+v \n", err0)
			log.Printf("error waiting for data copy in setdemagfield 1: %+v \n", err1)
			log.Printf("error waiting for data copy in setdemagfield 2: %+v \n", err2)
		}
		if NoDemagSpins.isZero() {
			if EnableNewellDemag && EnableDemag {
				log.Fatal("Cannot enable both Newell and brute force demag! \n")
			} else {
				if EnableNewellDemag {
					// Normal demag (Newell formulation), everywhere
					newellDemagConv().Exec(dst, M.Buffer(), geometry.Gpu(), msat, seqQueue, nil)
				}
				if EnableDemag {
					// Normal demag, everywhere
					demagConv().Exec(dst, M.Buffer(), geometry.Gpu(), msat, seqQueue, nil)
				}
			}
		} else {
			setMaskedDemagField(dst, msat, seqQueue, nil)
		}
	} else {
		opencl.Zero(dst) // will ADD other terms to it
		if err := seqQueue.Finish(); err != nil {
			log.Printf("error waiting for zero to finish in setdemagfield: %+v \n", err)
		}
	}
}

// Sets dst to the demag field, but cells where NoDemagSpins != 0 do not generate nor recieve field.
func setMaskedDemagField(dst *data.Slice, msat opencl.MSlice, q *cl.CommandQueue, events []*cl.Event) {
	// No-demag spins: mask-out geometry with zeros where NoDemagSpins is set,
	// so these spins do not generate a field

	buf := opencl.Buffer(SCALAR, geometry.Gpu().Size()) // masked-out geometry
	defer opencl.Recycle(buf)

	// obtain a copy of the geometry mask, which we can overwrite
	geom, r := geometry.Slice()
	if r {
		defer opencl.Recycle(geom)
	}
	data.Copy(buf, geom)
	if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
		log.Printf("error waiting for data copy in setmaskeddemagfield 0: %+v \n", err0)
		log.Printf("error waiting for data copy in setmaskeddemagfield 1: %+v \n", err1)
		log.Printf("error waiting for data copy in setmaskeddemagfield 2: %+v \n", err2)
	}

	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}

	// mask-out
	opencl.ZeroMask(buf, NoDemagSpins.gpuLUT1(), regions.Gpu(), queues, nil)

	// sync queues with sequential queue
	opencl.SyncQueues([]*cl.CommandQueue{q}, queues)

	// convolution with masked-out cells.
	if EnableDemag {
		// Normal demag, everywhere
		demagConv().Exec(dst, M.Buffer(), buf, msat, q, nil)
	}
	if EnableNewellDemag {
		// Normal demag (Newell formulation), everywhere
		newellDemagConv().Exec(dst, M.Buffer(), buf, msat, q, nil)
	}

	// sync queues with sequential queue
	opencl.SyncQueues(queues, []*cl.CommandQueue{q})

	// After convolution, mask-out the field in the NoDemagSpins cells
	// so they don't feel the field generated by others.
	opencl.ZeroMask(dst, NoDemagSpins.gpuLUT1(), regions.Gpu(), queues, nil)

	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		log.Printf("error waiting for all queues to finish after setmaskeddemagfield: %+v \n", err)
	}
}

// Sets dst to the full (unnormalized) magnetization in A/m
func SetMFull(dst *data.Slice) {
	// scale m by Msat...
	msat, rM := Msat.Slice()
	if rM {
		defer opencl.Recycle(msat)
	}
	// sync in the beginning
	if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
		log.Printf("error waiting for data copy in setmfull 0: %+v \n", err0)
		log.Printf("error waiting for data copy in setmfull 1: %+v \n", err1)
		log.Printf("error waiting for data copy in setmfull 2: %+v \n", err2)
	}
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}

	for c := 0; c < 3; c++ {
		opencl.Mul(dst.Comp(c), M.Buffer().Comp(c), msat, []*cl.CommandQueue{queues[c]}, nil)
	}

	// ...and by cell volume if applicable
	vol, rV := geometry.Slice()
	if rV {
		defer opencl.Recycle(vol)
	}
	if !vol.IsNil() {
		for c := 0; c < 3; c++ {
			opencl.Mul(dst.Comp(c), dst.Comp(c), vol, []*cl.CommandQueue{queues[c]}, nil)
		}
	}
	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		log.Printf("error waiting for all queues to finish after setmaskeddemagfield: %+v \n", err)
	}
}

// returns demag convolution, making sure it's initialized
func demagConv() *opencl.DemagConvolution {
	if conv_ == nil {
		SetBusy(true)
		defer SetBusy(false)
		kernel := mag.DemagKernel(Mesh().Size(), Mesh().PBC(), Mesh().CellSize(), DemagAccuracy, *Flag_cachedir)
		if err := opencl.WaitAllQueuesToFinish(); err != nil {
			log.Printf("error waiting all queues to finish in demagConv(): %+v \n", err)
		}
		conv_ = opencl.NewDemag(Mesh().Size(), Mesh().PBC(), kernel, *Flag_selftest, opencl.ClCmdQueue[0], nil)
	}
	return conv_
}

// returns demag convolution, making sure it's initialized
func newellDemagConv() *opencl.DemagConvolution {
	if conv_ == nil {
		SetBusy(true)
		defer SetBusy(false)
		seqQueue := opencl.ClCmdQueue[0]
		kernel := mag.NewellDemagKernel(Mesh().Size(), Mesh().PBC(), Mesh().CellSize(), asymptotic_radius, zero_self_demag, *Flag_cachedir)
		if err := opencl.WaitAllQueuesToFinish(); err != nil {
			log.Printf("error waiting all queues to finish in demagConv(): %+v \n", err)
		}
		conv_ = opencl.NewDemag(Mesh().Size(), Mesh().PBC(), kernel, *Flag_selftest, seqQueue, nil)
	}
	return conv_
}

// Returns the current demag energy in Joules.
func GetDemagEnergy() float64 {
	return -0.5 * cellVolume() * dot(&M_full, &B_demag)
}
