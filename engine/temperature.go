package engine

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	mag "github.com/seeder-research/uMagNUS/mag"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

var (
	Temp        = NewScalarParam("Temp", "K", "Temperature")
	E_therm     = NewScalarValue("E_therm", "J", "Thermal energy", GetThermalEnergy)
	Edens_therm = NewScalarField("Edens_therm", "J/m3", "Thermal energy density", AddThermalEnergyDensity)
	B_therm     thermField // Thermal effective field (T)
)

var AddThermalEnergyDensity = makeEdensAdder(&B_therm, -1)

// thermField calculates and caches thermal noise.
type thermField struct {
	seed      uint64            // seed for generator
	generator *opencl.Generator //
	noise     *data.Slice       // noise buffer
	step      int               // solver step corresponding to noise
	dt        float64           // solver timestep corresponding to noise
}

func init() {
	DeclFunc("ThermSeed", ThermSeed, "Set a random seed for thermal noise")
	registerEnergy(GetThermalEnergy, AddThermalEnergyDensity)
	B_therm.step = -1 // invalidate noise cache
	DeclROnly("B_therm", &B_therm, "Thermal field (T)")
}

func initRNG() uint64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Uint64()
}

func (b *thermField) UpdateSeed(seedVal *uint64) {
	if b.generator == nil {
		b.generator = opencl.NewGenerator("threefry")
	}
	if seedVal == nil {
		b.seed = initRNG()
	} else {
		b.seed = *seedVal
	}
	seqQueue := opencl.ClCmdQueue[0]
	b.generator.Init(&b.seed, seqQueue, nil)
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue at the end of thermfield.updateseed: %+v \n", err)
	}
}

func (b *thermField) AddTo(dst *data.Slice) {
	if !Temp.isZero() {
		// sync in the beginning
		if err := opencl.WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("error waiting for all queues to finish in thermfield.addto: %+v \n", err)
		}
		b.update()

		// checkout queues and execute kernel
		q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
		defer opencl.CheckinQueue(q1idx)
		defer opencl.CheckinQueue(q2idx)
		defer opencl.CheckinQueue(q3idx)
		queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
		seqQueue := opencl.ClCmdQueue[0]

		opencl.Add(dst, dst, b.noise, queues, nil)

		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
		if err := seqQueue.Finish(); err != nil {
			fmt.Printf("error waiting for queue to finish after thermfield.addto: %+v \n", err)
		}
	}
}

func (b *thermField) update() {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting all queues to finish in thermfield.update: %+v \n", err)
	}

	// we need to fix the time step here because solver will not yet have done it before the first step.
	// FixDt as an lvalue that sets Dt_si on change might be cleaner.
	if FixDt != 0 {
		Dt_si = FixDt
	}

	if b.generator == nil {
		b.generator = opencl.NewGenerator("threefry")
		b.seed = initRNG()
		b.UpdateSeed(&b.seed)
	}
	if b.noise == nil {
		b.noise = opencl.NewSlice(b.NComp(), b.Mesh().Size())
		// when noise was (re-)allocated it's invalid for sure.
		B_therm.step = -1
		B_therm.dt = -1
	}

	if Temp.isZero() {
		opencl.Memset(b.noise, 0, 0, 0)
		if err := opencl.H2DQueue.Finish(); err != nil {
			fmt.Printf("error waiting for queue if temp is zero in thermfield.update: %+v \n", err)
		}
		b.step = NSteps
		b.dt = Dt_si
		return
	}

	// keep constant during time step
	if NSteps == b.step && Dt_si == b.dt {
		return
	}

	// checkout queues and execute kernel
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
	seqQueue := opencl.ClCmdQueue[0]

	// after a bad step the timestep is rescaled and the noise should be rescaled accordingly, instead of redrawing the random numbers
	if NSteps == b.step && Dt_si != b.dt {
		for c := 0; c < 3; c++ {
			opencl.Madd2(b.noise.Comp(c), b.noise.Comp(c), b.noise.Comp(c), float32(math.Sqrt(b.dt/Dt_si)), 0., []*cl.CommandQueue{queues[c]}, nil)
		}
		b.dt = Dt_si
		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
		if err := seqQueue.Finish(); err != nil {
			fmt.Printf("error waiting for madd2 queues in thermfield.update: %+v \n", err)
		}
		return
	}

	if FixDt == 0 {
		Refer("leliaert2017")
		//uncomment to not allow adaptive step
		//util.Fatal("Finite temperature requires fixed time step. Set FixDt != 0.")
	}

	N := Mesh().NCell()
	k2_VgammaDt := 2 * mag.Kb / (GammaLL * cellVolume() * Dt_si)
	noise := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(noise)

	dst := b.noise
	ms := Msat.MSlice()
	defer ms.Recycle()
	temp := Temp.MSlice()
	defer temp.Recycle()
	alpha := Alpha.MSlice()
	defer alpha.Recycle()
	for i := 0; i < 3; i++ {
		b.generator.Normal(noise.DevPtr(0), int(N), queues[i], nil)
		opencl.SetTemperature(dst.Comp(i), noise, k2_VgammaDt, ms, temp, alpha, queues[i], nil)
	}
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue after thermfield.update: %+v \n", err)
	}
	b.step = NSteps
	b.dt = Dt_si
}

func GetThermalEnergy() float64 {
	if Temp.isZero() || relaxing {
		return 0
	} else {
		return -cellVolume() * dot(&M_full, &B_therm)
	}
}

// Seeds the thermal noise generator
func ThermSeed(seed int) {
	seedVal := (uint64)(seed)
	B_therm.UpdateSeed(&seedVal)
}

func (b *thermField) Mesh() *data.Mesh       { return Mesh() }
func (b *thermField) NComp() int             { return 3 }
func (b *thermField) Name() string           { return "Thermal field" }
func (b *thermField) Unit() string           { return "T" }
func (b *thermField) average() []float64     { return qAverageUniverse(b) }
func (b *thermField) EvalTo(dst *data.Slice) { EvalTo(b, dst) }
func (b *thermField) Slice() (*data.Slice, bool) {
	b.update()
	return b.noise, false
}
