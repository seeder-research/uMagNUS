package engine

import (
	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/mag"
	"github.com/seeder-research/uMagNUS/opencl"
	//"github.com/seeder-research/uMagNUS/opencl/cl"
	//"github.com/seeder-research/uMagNUS/util"
	//"fmt"
	"math"
	"math/rand"
	"time"
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
		b.generator = opencl.NewGenerator("mtgp")
	}
	if seedVal == nil {
		b.seed = initRNG()
	} else {
		b.seed = *seedVal
	}
	b.generator.Init(&b.seed, nil)
}

func (b *thermField) AddTo(dst *data.Slice) {
	if !Temp.isZero() {
		b.update()
		opencl.Add(dst, dst, b.noise)
	}
}

func (b *thermField) update() {
	// we need to fix the time step here because solver will not yet have done it before the first step.
	// FixDt as an lvalue that sets Dt_si on change might be cleaner.
	if FixDt != 0 {
		Dt_si = FixDt
	}

	if b.generator == nil {
		b.generator = opencl.NewGenerator("mtgp")
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
		b.step = NSteps
		b.dt = Dt_si
		return
	}

	// keep constant during time step
	if NSteps == b.step && Dt_si == b.dt {
		return
	}

	// after a bad step the timestep is rescaled and the noise should be rescaled accordingly, instead of redrawing the random numbers
	if NSteps == b.step && Dt_si != b.dt {
		for c := 0; c < 3; c++ {
			opencl.Madd2(b.noise.Comp(c), b.noise.Comp(c), b.noise.Comp(c), float32(math.Sqrt(b.dt/Dt_si)), 0.)
		}
		b.dt = Dt_si
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
		b.generator.Normal(noise.DevPtr(0), int(N), nil)
		opencl.SetTemperature(dst.Comp(i), noise, k2_VgammaDt, ms, temp, alpha)
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
