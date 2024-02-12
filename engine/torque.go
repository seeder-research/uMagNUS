package engine

import (
	"fmt"
	"reflect"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
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

	// For constant voltage type simulations (first interface)
	DisableVoltageInt1 = true
	ToMulFactorInt1    = true
	Vint1              = NewExcitation("Vint1", "", "Voltage applied to generate electrical current for interface 1 (same sign as J)")
	A1int1             = NewExcitation("A1int1", "", "First scale factor for calculating J from applied voltage for interface 1")
	A2int1             = NewExcitation("A2int1", "", "Second scale factor for calculating J from applied voltage for interface 1")
	VcmaCoeff1         = NewScalarParam("VcmaCoeff1", "J/m3/V", "voltage-controlled anisotropy constant for interface 1")
	AnisVCMAU1         = NewVectorParam("anisVCMAU1", "", "Voltage-controlled magnetic anisotropy direction for interface 1")
	// For constant voltage type simulations (second interface)
	DisableVoltageInt2 = true
	ToMulFactorInt2    = true
	Vint2              = NewExcitation("Vint2", "", "Voltage applied to generate electrical current for interface 2 (same sign as J)")
	A1int2             = NewExcitation("A1int2", "", "First scale factor for calculating J from applied voltage for interface 2")
	A2int2             = NewExcitation("A2int2", "", "Second scale factor for calculating J from applied voltage for interface 2")
	VcmaCoeff2         = NewScalarParam("VcmaCoeff2", "J/m3/V", "voltage-controlled anisotropy constant for interface 2")
	AnisVCMAU2         = NewVectorParam("anisVCMAU2", "", "Voltage-controlled magnetic anisotropy direction for interface 2")

	customTorques []Quantity // vector
)

func init() {
	Pol.setUniform([]float64{1}) // default spin polarization
	Lambda.Set(1)                // sensible default value (?).
	DeclFunc("AddTorqueTerm", AddTorqueTerm, "Adds an expression to total torque.")
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
	DeclVar("DisableVoltageInt1", &DisableVoltageInt1, "Disables voltage based calculation of Slonczewski torque through interface 1 (default=true)")
	DeclVar("ToMulFactorInt1", &ToMulFactorInt1, "Sets function converting voltage to current as multiply at interface 1 (default=true, divide if false)")
	DeclVar("DisableVoltageInt2", &DisableVoltageInt2, "Disables voltage based calculation of Slonczewski torque through interface 2 (default=true)")
	DeclVar("ToMulFactorInt2", &ToMulFactorInt2, "Sets function converting voltage to current as multiply at interface 2 (default=true, divide if false)")
	DeclFunc("RemoveCustomTorques", RemoveCustomTorques, "Removes all custom torques again")
}

// Removes all customfields
func RemoveCustomTorques() {
	customTorques = nil
}

// AddFieldTerm adds an effective field function (returning Teslas) to B_eff.
// Be sure to also add the corresponding energy term using AddEnergyTerm.
func AddTorqueTerm(b Quantity) {
	customTorques = append(customTorques, b)
}

// AddCustomTorque evaluates the user-defined custom torque terms
// and adds the result to dst.
func AddCustomTorques(dst *data.Slice) {
	if len(customTorques) > 0 {
		// sync in the beginning
		if err := opencl.WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("error waiting for queues to finish before addcustomtorques: %+v \n", err)
		}
		// checkout queues for execution
		q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
		defer opencl.CheckinQueue(q1idx)
		defer opencl.CheckinQueue(q2idx)
		defer opencl.CheckinQueue(q3idx)

		for _, term := range customTorques {
			buf := ValueOf(term)
			queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
			opencl.Add(dst, dst, buf, queues, nil)

			// sync to prevent buf from being recycled before the kernel finishes
			if err := opencl.WaitAllQueuesToFinish(); err != nil {
				fmt.Printf("error waiting for queues to finish in addcustomtorques: %+v \n", err)
			}

			opencl.Recycle(buf)
		}
	}
}

// Sets dst to the current total torque
func SetTorque(dst *data.Slice) {

	// sync in the beginning
	if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
		fmt.Printf("error waiting for queues to finish in settorque 0: %+v \n", err0)
		fmt.Printf("error waiting for queues to finish in settorque 1: %+v \n", err1)
		fmt.Printf("error waiting for queues to finish in settorque 2: %+v \n", err2)
	}

	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	seqQueue := opencl.ClCmdQueue[0]

	// SetLLTorque(dst) - torque from fields ------------
	SetEffectiveField(dst) // calc and store B_eff
	alpha := Alpha.MSlice()
	defer alpha.Recycle()

	if Precess {
		opencl.LLTorque(dst, M.Buffer(), dst, alpha, opencl.ClCmdQueue[q1idx], nil) // overwrite dst with torque
	} else {
		opencl.LLNoPrecess(dst, M.Buffer(), dst, opencl.ClCmdQueue[q1idx], nil)
	}
	// sync queue with seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{opencl.ClCmdQueue[q1idx]})

	// Checkout buffers for remaining torques
	var sttBuf1, sttBuf2, sttBuf3, sttBuf4 *data.Slice

	msat := Msat.MSlice()
	defer msat.Recycle()

	// AddSTTorque(dst) - stt torque -----------
	if J.isZero() == false {
		j := J.MSlice()
		defer j.Recycle()
		xi := Xi.MSlice()
		defer xi.Recycle()
		pol := Pol.MSlice()
		defer pol.Recycle()
		fixedP := FixedLayer.MSlice()
		defer fixedP.Recycle()
		lambda := Lambda.MSlice()
		defer lambda.Recycle()
		epsPrime := EpsilonPrime.MSlice()
		defer epsPrime.Recycle()
		thickness := FreeLayerThickness.MSlice()
		defer thickness.Recycle()

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
			// checkout buffer and queue for execution
			sttBuf1 = opencl.Buffer(dst.NComp(), dst.Size())
			defer opencl.Recycle(sttBuf1)
			opencl.AddZhangLiTorque(sttBuf1, M.Buffer(), msat, j, alpha, xi, pol, Mesh(), opencl.ClCmdQueue[q2idx], nil)
			opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{opencl.ClCmdQueue[q2idx]})
		}
		if !DisableSlonczewskiTorque && !FixedLayer.isZero() {
			// checkout buffer and queue for execution
			sttBuf2 = opencl.Buffer(dst.NComp(), dst.Size())
			defer opencl.Recycle(sttBuf2)
			opencl.AddSlonczewskiTorque2(sttBuf2, M.Buffer(),
				msat, j, fixedP, alpha, pol, lambda, epsPrime,
				thickness,
				CurrentSignFromFixedLayerPosition[fixedLayerPosition],
				Mesh(),
				opencl.ClCmdQueue[q3idx], nil)
			// sync queue with seqQueue
			opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{opencl.ClCmdQueue[q3idx]})
		}
		opencl.SyncQueues([]*cl.CommandQueue{opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}, []*cl.CommandQueue{seqQueue})
		queues := []*cl.CommandQueue{seqQueue, opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
		if !DisableZhangLiTorque {
			if !DisableSlonczewskiTorque && !FixedLayer.isZero() {
				opencl.Madd3(dst, dst, sttBuf1, sttBuf2, 1., 1., 1., queues, nil)
			} else {
				opencl.Add(dst, dst, sttBuf1, queues, nil)
			}
		} else {
			if !DisableSlonczewskiTorque && !FixedLayer.isZero() {
				opencl.Add(dst, dst, sttBuf2, queues, nil)
			}
		}
	}

	// AddSTTorque1(dst) - stt torque ----------------
	if (Jint1.isZero() && Vint1.isZero()) == false {
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
			// Checkout buffer and queue for executing kernel
			sttBuf3 = opencl.Buffer(dst.NComp(), dst.Size())
			defer opencl.Recycle(sttBuf3)
			q4idx := opencl.CheckoutQueue()
			defer opencl.CheckinQueue(q4idx)

			vapp, rec := Vint1.Slice()
			if rec {
				defer opencl.Recycle(vapp)
			}
			Jcurr := opencl.Buffer(vapp.NComp(), vapp.Size())
			defer opencl.Recycle(Jcurr)
			if !DisableVoltageInt1 {
				a1int, rec := A1int1.Slice()
				if rec {
					defer opencl.Recycle(a1int)
				}
				a2int, rec := A2int1.Slice()
				if rec {
					defer opencl.Recycle(a2int)
				}
				JfromV(vapp, a1int, a2int, M.Buffer(), fl, Jcurr, ToMulFactorInt1)
				Jint1.AddTo(Jcurr)
			}
			j := opencl.ToMSlice(Jcurr)
			defer j.Recycle()
			fixedP := FixedLayer1.MSlice()
			defer fixedP.Recycle()
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
			opencl.AddOommfSlonczewskiTorque(sttBuf3, M.Buffer(),
				msat, j, fixedP, alpha, pfix, pfree, lambdafix, lambdafree, epsPrime, Mesh(), opencl.ClCmdQueue[q4idx], nil)
			// sync queues and merge result
			opencl.SyncQueues([]*cl.CommandQueue{seqQueue, opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}, []*cl.CommandQueue{opencl.ClCmdQueue[q4idx]})
			opencl.Add(dst, dst, sttBuf3, []*cl.CommandQueue{seqQueue, opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx]}, nil)
			// sync queues
			opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx]})
		}
	}

	// AddSTTorque2(dst) - stt torque ----------------
	if (Jint2.isZero() && Vint2.isZero()) == false {
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
			// Checkout buffer and queue for executing kernel
			sttBuf4 = opencl.Buffer(dst.NComp(), dst.Size())
			defer opencl.Recycle(sttBuf4)
			q4idx := opencl.CheckoutQueue()
			defer opencl.CheckinQueue(q4idx)

			vapp, rec := Vint2.Slice()
			if rec {
				defer opencl.Recycle(vapp)
			}
			Jcurr := opencl.Buffer(vapp.NComp(), vapp.Size())
			defer opencl.Recycle(Jcurr)
			if !DisableVoltageInt2 {
				a1int, rec := A1int2.Slice()
				if rec {
					defer opencl.Recycle(a1int)
				}
				a2int, rec := A2int2.Slice()
				if rec {
					defer opencl.Recycle(a2int)
				}
				JfromV(vapp, a1int, a2int, M.Buffer(), fl, Jcurr, ToMulFactorInt2)
				Jint2.AddTo(Jcurr)
			}
			j := opencl.ToMSlice(Jcurr)
			defer j.Recycle()
			fixedP := FixedLayer2.MSlice()
			defer fixedP.Recycle()
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
			opencl.AddOommfSlonczewskiTorque(sttBuf4, M.Buffer(),
				msat, j, fixedP, alpha, pfix, pfree, lambdafix, lambdafree, epsPrime, Mesh(), opencl.ClCmdQueue[q4idx], nil)
			// sync queues and merge result
			opencl.SyncQueues([]*cl.CommandQueue{seqQueue, opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}, []*cl.CommandQueue{opencl.ClCmdQueue[q4idx]})
			opencl.Add(dst, dst, sttBuf3, []*cl.CommandQueue{seqQueue, opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx]}, nil)
			// sync queues
			opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx]})
		}
	}
	// -----

	AddRegionLinkSpinTorque(dst) // special stt torque (synchronous)
	AddCustomTorques(dst)        // custom torques (synchronous)
	FreezeSpins(dst)             // zero the torques at locations where spins are frozen.
}

// Sets dst to the current Landau-Lifshitz torque
func SetLLTorque(dst *data.Slice) {
	SetEffectiveField(dst) // calc and store B_eff
	alpha := Alpha.MSlice()
	defer alpha.Recycle()
	seqQueue := opencl.ClCmdQueue[0]
	if Precess {
		opencl.LLTorque(dst, M.Buffer(), dst, alpha, seqQueue, nil) // overwrite dst with torque
	} else {
		opencl.LLNoPrecess(dst, M.Buffer(), dst, seqQueue, nil)
	}
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue after setlltorque: %+v \n", err)
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
		opencl.AddZhangLiTorque(dst, M.Buffer(), msat, j, alpha, xi, pol, Mesh(), opencl.ClCmdQueue[0], nil)
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
		// Checkout a queue and execute synchronously
		opencl.AddSlonczewskiTorque2(dst, M.Buffer(),
			msat, j, fixedP, alpha, pol, lambda, epsPrime,
			thickness,
			CurrentSignFromFixedLayerPosition[fixedLayerPosition],
			Mesh(), opencl.ClCmdQueue[0], nil)
	}
}

// Adds the current spin transfer torque from first additional source to dst
func AddSTTorque1(dst *data.Slice) {
	if Jint1.isZero() && Vint1.isZero() {
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
		vapp, rec := Vint1.Slice()
		if rec {
			defer opencl.Recycle(vapp)
		}
		Jcurr := opencl.Buffer(vapp.NComp(), vapp.Size())
		defer opencl.Recycle(Jcurr)
		if !DisableVoltageInt1 {
			a1int, rec := A1int1.Slice()
			if rec {
				defer opencl.Recycle(a1int)
			}
			a2int, rec := A2int1.Slice()
			if rec {
				defer opencl.Recycle(a2int)
			}
			JfromV(vapp, a1int, a2int, M.Buffer(), fl, Jcurr, ToMulFactorInt1)
			Jint1.AddTo(Jcurr)
		}
		j := opencl.ToMSlice(Jcurr)
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
			msat, j, fixedP, alpha, pfix, pfree, lambdafix, lambdafree, epsPrime, Mesh(), opencl.ClCmdQueue[0], nil)
	}
}

// Adds the current spin transfer torque from second additional source to dst
func AddSTTorque2(dst *data.Slice) {
	if Jint2.isZero() && Vint2.isZero() {
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
		vapp, rec := Vint2.Slice()
		if rec {
			defer opencl.Recycle(vapp)
		}
		Jcurr := opencl.Buffer(vapp.NComp(), vapp.Size())
		defer opencl.Recycle(Jcurr)
		if !DisableVoltageInt2 {
			a1int, rec := A1int2.Slice()
			if rec {
				defer opencl.Recycle(a1int)
			}
			a2int, rec := A2int2.Slice()
			if rec {
				defer opencl.Recycle(a2int)
			}
			JfromV(vapp, a1int, a2int, M.Buffer(), fl, Jcurr, ToMulFactorInt2)
			Jint2.AddTo(Jcurr)
		}
		j := opencl.ToMSlice(Jcurr)
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
			msat, j, fixedP, alpha, pfix, pfree, lambdafix, lambdafree, epsPrime, Mesh(), opencl.ClCmdQueue[0], nil)
	}
}

func JfromV(Vappl, A1, A2, m, refM, Jcurr *data.Slice, ToMulFactor bool) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues in jfromv: %+v \n", err)
	}
	cellSz := M.Mesh().CellSize()
	xSz, ySz, zSz := cellSz[X], cellSz[Y], cellSz[Z]
	xArea := make([]float64, 3)
	xArea[X], xArea[Y], xArea[Z] = ySz*zSz, xSz*zSz, xSz*ySz

	a1int := opencl.Buffer(1, A1.Size())
	a2int := opencl.Buffer(1, A1.Size())
	factor := opencl.Buffer(1, A1.Size())
	factor1 := opencl.Buffer(1, A1.Size())
	dp := opencl.Buffer(1, A1.Size())
	defer opencl.Recycle(a1int)
	defer opencl.Recycle(a2int)
	defer opencl.Recycle(factor)
	defer opencl.Recycle(factor1)
	defer opencl.Recycle(dp)
	opencl.Zero(dp)

	// checkout queues for execution
	seqQueue := opencl.ClCmdQueue[0]
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	q4idx, q5idx, q6idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	q7idx := opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	defer opencl.CheckinQueue(q4idx)
	defer opencl.CheckinQueue(q5idx)
	defer opencl.CheckinQueue(q6idx)
	defer opencl.CheckinQueue(q7idx)
	queues1 := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
	queues2 := []*cl.CommandQueue{opencl.ClCmdQueue[q4idx], opencl.ClCmdQueue[q5idx], opencl.ClCmdQueue[q6idx]}
	// sync queues
	opencl.SyncQueues([]*cl.CommandQueue{opencl.ClCmdQueue[q7idx]}, []*cl.CommandQueue{seqQueue}) // dp

	opencl.AddDotProduct(dp, float32(1.0), m, refM, opencl.ClCmdQueue[q7idx], nil)
	// sync queues
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{opencl.ClCmdQueue[q7idx]})

	for ii := 0; ii < A1.NComp(); ii++ {
		opencl.Madd2(a1int, A1.Comp(ii), A2.Comp(ii), float32(0.5), float32(0.5), []*cl.CommandQueue{queues1[ii]}, nil)
		opencl.Madd2(a2int, A1.Comp(ii), A2.Comp(ii), float32(0.5), float32(-0.5), []*cl.CommandQueue{queues2[ii]}, nil)

		// sync queues
		opencl.SyncQueues([]*cl.CommandQueue{queues2[ii]}, []*cl.CommandQueue{opencl.ClCmdQueue[q7idx]})
		opencl.Mul(factor1, a2int, dp, []*cl.CommandQueue{queues2[ii]}, nil)
		// sync queues
		opencl.SyncQueues([]*cl.CommandQueue{queues2[ii]}, []*cl.CommandQueue{queues1[ii]})

		if ToMulFactor {
			opencl.Madd2(factor, a1int, factor1, float32(float64(1.0)/xArea[ii]), float32(float64(1.0)/xArea[ii]), []*cl.CommandQueue{queues2[ii]}, nil)
			opencl.Mul(Jcurr.Comp(ii), Vappl.Comp(ii), factor, []*cl.CommandQueue{queues2[ii]}, nil)
		} else {
			opencl.Madd2(factor, a1int, factor1, float32(xArea[ii]), float32(xArea[ii]), []*cl.CommandQueue{queues2[ii]}, nil)
			opencl.Div(Jcurr.Comp(ii), Vappl.Comp(ii), factor, []*cl.CommandQueue{queues2[ii]}, nil)
		}
	}

	// Sync all queues
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish after jfromv: %+v \n", err)
	}
}

func FreezeSpins(dst *data.Slice) {
	if !FrozenSpins.isZero() {
		// sync in the beginning
		if err := opencl.WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("error waiting for queues to finish in freezespins: %+v \n", err)
		}

		// Checkout queues and execute kernel
		q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
		defer opencl.CheckinQueue(q1idx)
		defer opencl.CheckinQueue(q2idx)
		defer opencl.CheckinQueue(q3idx)
		queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
		opencl.ZeroMask(dst, FrozenSpins.gpuLUT1(), regions.Gpu(), queues, nil)

		// sync before returning
		if err := opencl.WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("error waiting for queues to finish after freezespins: %+v \n", err)
		}
	}
}

func GetMaxTorque() float64 {
	torque := ValueOf(Torque)
	defer opencl.Recycle(torque)
	seqQueue := opencl.ClCmdQueue[0]
	return opencl.MaxVecNorm(torque, seqQueue, nil)
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
