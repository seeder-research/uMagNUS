package engine

import (
	"reflect"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl"
	"github.com/seeder-research/uMagNUS/util"
)

var (
	Alpha                            = NewScalarParam("alpha", "", "Landau-Lifshitz damping constant")
	Xi                               = NewScalarParam("xi", "", "Non-adiabaticity of spin-transfer-torque")
	Pol                              = NewScalarParam("Pol", "", "Electrical current polarization")
	Lambda                           = NewScalarParam("Lambda", "", "Slonczewski Λ parameter")
	EpsilonPrime                     = NewScalarParam("EpsilonPrime", "", "Slonczewski secondairy STT term ε'")
	FrozenSpins                      = NewScalarParam("frozenspins", "", "Defines spins that should be fixed") // 1 - frozen, 0 - free. TODO: check if it only contains 0/1 values
	FreeLayerThickness               = NewScalarParam("FreeLayerThickness", "m", "Slonczewski free layer thickness (if set to zero (default), then the thickness will be deduced from the mesh size)")
	FixedLayer                       = NewExcitation("FixedLayer", "", "Slonczewski fixed layer polarization")
	Torque                           = NewVectorField("torque", "T", "Total torque/γ0", SetTorque)
	LLTorque                         = NewVectorField("LLtorque", "T", "Landau-Lifshitz torque/γ0", SetLLTorque)
	STTorque                         = NewVectorField("STTorque", "T", "Spin-transfer torque/γ0", AddSTTorque)
	STTorqueInt1                     = NewVectorField("STTorqueInt1", "T", "Spin-transfer torque/γ0", AddSTTorque1)
	STTorqueInt2                     = NewVectorField("STTorqueInt2", "T", "Spin-transfer torque/γ0", AddSTTorque2)
	J                                = NewExcitation("J", "A/m2", "Electrical current density")
	MaxTorque                        = NewScalarValue("maxTorque", "T", "Maximum torque/γ0, over all cells", GetMaxTorque)
	GammaLL                  float64 = 1.7595e11 // Gyromagnetic ratio of spins, in rad/Ts
	Precess                          = true
	DisableZhangLiTorque             = false
	DisableSlonczewskiTorque         = false
	fixedLayerPosition               = FIXEDLAYER_TOP // instructs uMagNUS how free and fixed layers are stacked along +z direction

	// For first additional source of spin torque
	DisableSlonczewskiTorque1 = true
	Pfree1                    = NewScalarParam("Pfree1", "", "Electrical current polarization (free layer side) for interface 1")
	Pfixed1                   = NewScalarParam("Pfixed1", "", "Electrical current polarization (fixed layer side) for interface 1")
	Lambdafree1               = NewScalarParam("Lambdafree1", "", "Slonczewski Λ_free parameter for interface 1")
	Lambdafixed1              = NewScalarParam("Lambdafixed1", "", "Slonczewski Λ_fixed parameter for interface 1")
	EpsilonPrime1             = NewScalarParam("EpsilonPrime1", "", "Slonczewski secondairy STT term ε' for interface 1")
	FixedLayer1               = NewExcitation("FixedLayer1", "", "Slonczewski fixed layer polarization for interface 1")
	Jint1                     = NewExcitation("Jint1", "A/m2", "Electrical current density through interface 1")
	// For second additional source of spin torque
	DisableSlonczewskiTorque2 = true
	Pfree2                    = NewScalarParam("Pfree2", "", "Electrical current polarization (free layer side) for interface 2")
	Pfixed2                   = NewScalarParam("Pfixed2", "", "Electrical current polarization (fixed layer side) for interface 2")
	Lambdafree2               = NewScalarParam("Lambdafree2", "", "Slonczewski Λ_free parameter for interface 2")
	Lambdafixed2              = NewScalarParam("Lambdafixed2", "", "Slonczewski Λ_fixed parameter for interface 2")
	EpsilonPrime2             = NewScalarParam("EpsilonPrime2", "", "Slonczewski secondairy STT term ε' for interface 2")
	FixedLayer2               = NewExcitation("FixedLayer2", "", "Slonczewski fixed layer polarization for interface 2")
	Jint2                     = NewExcitation("Jint2", "A/m2", "Electrical current density through interface 2")
)

func init() {
	Pol.setUniform([]float64{1}) // default spin polarization
	Lambda.Set(1)                // sensible default value (?).
	DeclVar("GammaLL", &GammaLL, "Gyromagnetic ratio in rad/Ts")
	DeclVar("DisableZhangLiTorque", &DisableZhangLiTorque, "Disables Zhang-Li torque (default=false)")
	DeclVar("DisableSlonczewskiTorque", &DisableSlonczewskiTorque, "Disables Slonczewski torque (default=false)")
	DeclVar("DoPrecess", &Precess, "Enables LL precession (default=true)")
	DeclLValue("FixedLayerPosition", &flposition{}, "Position of the fixed layer: FIXEDLAYER_TOP, FIXEDLAYER_BOTTOM (default=FIXEDLAYER_TOP)")
	DeclROnly("FIXEDLAYER_TOP", FIXEDLAYER_TOP, "FixedLayerPosition = FIXEDLAYER_TOP instructs uMagNUS that fixed layer is on top of the free layer")
	DeclROnly("FIXEDLAYER_BOTTOM", FIXEDLAYER_BOTTOM, "FixedLayerPosition = FIXEDLAYER_BOTTOM instructs uMagNUS that fixed layer is underneath of the free layer")

	// For two OOMMF type Slonczewski interfaces
	Pfree1.setUniform([]float64{1})  // default spin polarization
	Pfixed1.setUniform([]float64{1}) // default spin polarization
	Lambdafree1.Set(1)               // sensible default value (?). TODO: should not be zero
	Lambdafixed1.Set(1)              // sensible default value (?). TODO: should not be zero
	Pfree2.setUniform([]float64{1})  // default spin polarization
	Pfixed2.setUniform([]float64{1}) // default spin polarization
	Lambdafree2.Set(1)               // sensible default value (?). TODO: should not be zero
	Lambdafixed2.Set(1)              // sensible default value (?). TODO: should not be zero
	DeclVar("DisableSlonczewskiTorque1", &DisableSlonczewskiTorque1, "Disables Slonczewski torque through interface 1 (default=true)")
	DeclVar("DisableSlonczewskiTorque2", &DisableSlonczewskiTorque2, "Disables Slonczewski torque through interface 2 (default=true)")
}

// Sets dst to the current total torque
func SetTorque(dst *data.Slice) {
	SetLLTorque(dst)
	AddSTTorque(dst)
	AddSTTorque1(dst)
	AddSTTorque2(dst)
	FreezeSpins(dst)
}

// Sets dst to the current Landau-Lifshitz torque
func SetLLTorque(dst *data.Slice) {
	SetEffectiveField(dst) // calc and store B_eff
	alpha := Alpha.MSlice()
	defer alpha.Recycle()
	if Precess {
		opencl.LLTorque(dst, M.Buffer(), dst, alpha) // overwrite dst with torque
	} else {
		opencl.LLNoPrecess(dst, M.Buffer(), dst)
	}
}

// Adds the current spin transfer torque to dst
func AddSTTorque(dst *data.Slice) {
	if J.isZero() {
		return
	}
	util.AssertMsg(!Pol.isZero(), "spin polarization should not be 0")
	jspin, rec := J.Slice()
	if rec {
		defer opencl.Recycle(jspin)
	}
	fl, rec := FixedLayer.Slice()
	if rec {
		defer opencl.Recycle(fl)
	}
	if !DisableZhangLiTorque {
		msat := Msat.MSlice()
		defer msat.Recycle()
		j := J.MSlice()
		defer j.Recycle()
		alpha := Alpha.MSlice()
		defer alpha.Recycle()
		xi := Xi.MSlice()
		defer xi.Recycle()
		pol := Pol.MSlice()
		defer pol.Recycle()
		opencl.AddZhangLiTorque(dst, M.Buffer(), msat, j, alpha, xi, pol, Mesh())
	}
	if !DisableSlonczewskiTorque && !FixedLayer.isZero() {
		msat := Msat.MSlice()
		defer msat.Recycle()
		j := J.MSlice()
		defer j.Recycle()
		fixedP := FixedLayer.MSlice()
		defer fixedP.Recycle()
		alpha := Alpha.MSlice()
		defer alpha.Recycle()
		pol := Pol.MSlice()
		defer pol.Recycle()
		lambda := Lambda.MSlice()
		defer lambda.Recycle()
		epsPrime := EpsilonPrime.MSlice()
		defer epsPrime.Recycle()
		thickness := FreeLayerThickness.MSlice()
		defer thickness.Recycle()
		opencl.AddSlonczewskiTorque2(dst, M.Buffer(),
			msat, j, fixedP, alpha, pol, lambda, epsPrime,
			thickness,
			CurrentSignFromFixedLayerPosition[fixedLayerPosition],
			Mesh())
	}
}

// Adds the current spin transfer torque from first additional source to dst
func AddSTTorque1(dst *data.Slice) {
	if Jint1.isZero() {
		return
	}
	util.AssertMsg(!Pfree1.isZero(), "spin polarization (Pfree1) should not be 0")
	util.AssertMsg(!Pfixed1.isZero(), "spin polarization (Pfixed1) should not be 0")
	jspin, rec := Jint1.Slice()
	if rec {
		defer opencl.Recycle(jspin)
	}
	fl, rec := FixedLayer1.Slice()
	if rec {
		defer opencl.Recycle(fl)
	}
	if !DisableSlonczewskiTorque1 && !FixedLayer1.isZero() {
		msat := Msat.MSlice()
		defer msat.Recycle()
		j := Jint1.MSlice()
		defer j.Recycle()
		fixedP := FixedLayer1.MSlice()
		defer fixedP.Recycle()
		alpha := Alpha.MSlice()
		defer alpha.Recycle()
		pfix := Pfixed1.MSlice()
		defer pfix.Recycle()
		pfree := Pfree1.MSlice()
		defer pfree.Recycle()
		lambdafree := Lambdafree1.MSlice()
		defer lambdafree.Recycle()
		lambdafix := Lambdafixed1.MSlice()
		defer lambdafix.Recycle()
		epsPrime := EpsilonPrime1.MSlice()
		defer epsPrime.Recycle()
		opencl.AddOommfSlonczewskiTorque(dst, M.Buffer(),
			msat, j, fixedP, alpha, pfix, pfree, lambdafix, lambdafree, epsPrime, Mesh())
	}
}

// Adds the current spin transfer torque from second additional source to dst
func AddSTTorque2(dst *data.Slice) {
	if Jint2.isZero() {
		return
	}
	util.AssertMsg(!Pfree2.isZero(), "spin polarization (Pfree1) should not be 0")
	util.AssertMsg(!Pfixed2.isZero(), "spin polarization (Pfixed1) should not be 0")
	jspin, rec := Jint2.Slice()
	if rec {
		defer opencl.Recycle(jspin)
	}
	fl, rec := FixedLayer2.Slice()
	if rec {
		defer opencl.Recycle(fl)
	}
	if !DisableSlonczewskiTorque2 && !FixedLayer2.isZero() {
		msat := Msat.MSlice()
		defer msat.Recycle()
		j := Jint2.MSlice()
		defer j.Recycle()
		fixedP := FixedLayer1.MSlice()
		defer fixedP.Recycle()
		alpha := Alpha.MSlice()
		defer alpha.Recycle()
		pfix := Pfixed2.MSlice()
		defer pfix.Recycle()
		pfree := Pfree2.MSlice()
		defer pfree.Recycle()
		lambdafree := Lambdafree2.MSlice()
		defer lambdafree.Recycle()
		lambdafix := Lambdafixed2.MSlice()
		defer lambdafix.Recycle()
		epsPrime := EpsilonPrime2.MSlice()
		defer epsPrime.Recycle()
		opencl.AddOommfSlonczewskiTorque(dst, M.Buffer(),
			msat, j, fixedP, alpha, pfix, pfree, lambdafix, lambdafree, epsPrime, Mesh())
	}
}

func FreezeSpins(dst *data.Slice) {
	if !FrozenSpins.isZero() {
		opencl.ZeroMask(dst, FrozenSpins.gpuLUT1(), regions.Gpu())
	}
}

func GetMaxTorque() float64 {
	torque := ValueOf(Torque)
	defer opencl.Recycle(torque)
	return opencl.MaxVecNorm(torque)
}

type FixedLayerPosition int

const (
	FIXEDLAYER_TOP FixedLayerPosition = iota + 1
	FIXEDLAYER_BOTTOM
)

var (
	CurrentSignFromFixedLayerPosition = map[FixedLayerPosition]float64{
		FIXEDLAYER_TOP:    1.0,
		FIXEDLAYER_BOTTOM: -1.0,
	}
)

type flposition struct{}

func (*flposition) Eval() interface{} { return fixedLayerPosition }
func (*flposition) SetValue(v interface{}) {
	drainOutput()
	fixedLayerPosition = v.(FixedLayerPosition)
}
func (*flposition) Type() reflect.Type { return reflect.TypeOf(FixedLayerPosition(FIXEDLAYER_TOP)) }
